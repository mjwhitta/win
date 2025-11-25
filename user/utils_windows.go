//go:build windows

package user

import (
	"bytes"
	"encoding/binary"
	"slices"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

func adjustToken(p *Privilege) (e error) {
	var t windows.Token
	var tp *windows.Tokenprivileges = &windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{
			{Attributes: p.Attributes, Luid: p.LUID},
		},
	}

	e = windows.OpenProcessToken(
		p.proc,
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&t,
	)
	if e != nil {
		return errors.Newf("failed to get process token: %w", e)
	}
	defer func() {
		if e == nil {
			e = t.Close()
		}
	}()

	e = windows.AdjustTokenPrivileges(
		t,
		false,
		tp,
		uint32(unsafe.Sizeof(*tp)),
		nil,
		nil,
	)
	if e != nil {
		return errors.Newf("failed to adjust token privileges: %w", e)
	}

	return nil
}

func getGroupAttrs(attributes uint32) []string {
	var attrs []string
	var keys []uint32
	var mask map[uint32]string = map[uint32]string{
		0x1:        "Mandatory group",
		0x2:        "Enabled by default",
		0x4:        "Enabled group",
		0x8:        "Group owner",
		0x10:       "Group used for deny only",
		0x20:       "Integrity",
		0x40:       "Integrity enabled",
		0xc0000000: "Logon ID",
		0x20000000: "Local Group", // Resource
	}
	var valid uint32 = 0xe000007f

	if attributes|valid != valid {
		return nil
	}

	for k := range mask {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, k := range keys {
		if attributes&k > 0 {
			attrs = append(attrs, mask[k])
		}
	}

	return attrs
}

func getGroupNameAndType(sid *windows.SID) (string, string) {
	var account string
	var domain string
	var name string
	var theType uint32
	var types map[uint32]string = map[uint32]string{
		0:  "Unknown SID type",
		1:  "User",
		2:  "Group",
		3:  "Domain",
		4:  "Alias",
		5:  "Well-known group",
		6:  "Deleted",
		7:  "Invalid",
		8:  "Computer",
		10: "Label",
	}

	account, domain, theType, _ = sid.LookupAccount(".")

	if account == "None" {
		return "", ""
	}

	if strings.HasPrefix(account, "LogonSessionId_") {
		return "", ""
	}

	name = domain + "\\" + account
	if domain == "" {
		name = account
	}

	if accType, ok := types[theType]; ok {
		return name, accType
	}

	return name, types[0]
}

func getPrivName(luid windows.LUID) (string, error) {
	var b []byte
	var e error
	var n int
	//nolint:mnd // Shift 32 bits left
	var l int64 = (int64(luid.HighPart) << 32) + int64(luid.LowPart)

	if e = w32.LookupPrivilegeName("", l, &b, &n); e != nil {
		b = make([]byte, n)
		if e = w32.LookupPrivilegeName("", l, &b, &n); e != nil {
			e = errors.Newf("failed to lookup privilege name: %w", e)
			return "", e
		}
	}

	return string(b), nil
}

func getPrivDesc(name string) (string, error) {
	var b []byte
	var e error
	var n int

	e = w32.LookupPrivilegeDisplayName("", name, &b, &n)
	if e != nil {
		b = make([]byte, n)

		e = w32.LookupPrivilegeDisplayName("", name, &b, &n)
		if e != nil {
			return "", errors.Newf(
				"failed to lookup privilege description: %w",
				e,
			)
		}
	}

	return string(b), nil
}

func output(section string, hdrs []string, data [][]string) string {
	var lines []string
	var sb strings.Builder
	var width []int = make([]int, len(hdrs))

	// Initial max width
	for i, col := range hdrs {
		width[i] = len(col)
	}

	// Find max width
	for _, row := range data {
		for i, col := range row {
			if len(col) > width[i] {
				width[i] = len(col)
			}
		}
	}

	// Section
	lines = append(lines, "")
	lines = append(lines, section)
	lines = append(lines, strings.Repeat("-", len(section)))
	lines = append(lines, "")

	// Headers
	for i, col := range hdrs {
		sb.WriteString(col)
		sb.WriteString(strings.Repeat(" ", width[i]-len(col)))
		sb.WriteString(" ")
	}

	lines = append(lines, sb.String())

	// Dividers
	sb.Reset()

	for i := range hdrs {
		sb.WriteString(strings.Repeat("=", width[i]) + " ")
	}

	lines = append(lines, sb.String())

	// Data
	for _, row := range data {
		sb.Reset()

		for i, col := range row {
			sb.WriteString(col)
			sb.WriteString(strings.Repeat(" ", width[i]-len(col)))
			sb.WriteString(" ")
		}

		lines = append(lines, sb.String())
	}

	// Print
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}

func privsFromBytes(
	b []byte,
	n uint32,
	proc windows.Handle,
) ([]*Privilege, error) {
	var attrs uint32
	var buf *bytes.Buffer = bytes.NewBuffer(b)
	var e error
	var privs []*Privilege

	if e = binary.Read(buf, binary.LittleEndian, &n); e != nil {
		e = errors.Newf("failed to read number of privileges: %w", e)
		return nil, e
	}

	privs = make([]*Privilege, n)

	for i := range privs {
		privs[i] = &Privilege{proc: proc}

		e = binary.Read(buf, binary.LittleEndian, &privs[i].LUID)
		if e != nil {
			return nil, errors.Newf("failed to read LUID: %w", e)
		}

		if privs[i].Name, e = getPrivName(privs[i].LUID); e != nil {
			return nil, e
		}

		privs[i].Description, e = getPrivDesc(privs[i].Name)
		if e != nil {
			return nil, e
		}

		e = binary.Read(buf, binary.LittleEndian, &attrs)
		if e != nil {
			e = errors.Newf("failed to read attributes: %w", e)
			return nil, e
		}

		privs[i].Attributes = attrs
	}

	return privs, nil
}

func procOrDefault(proc []windows.Handle) windows.Handle {
	if len(proc) == 0 {
		return windows.CurrentProcess()
	}

	return proc[0]
}

func tokenOrDefault(proc []windows.Handle) (windows.Token, error) {
	var e error
	var t windows.Token

	if len(proc) == 0 {
		return windows.GetCurrentProcessToken(), nil
	}

	e = windows.OpenProcessToken(proc[0], windows.TOKEN_QUERY, &t)
	if e != nil {
		e = errors.Newf(
			"failed to open process token for %v: %w",
			proc[0],
			e,
		)

		return 0, e
	}

	return t, nil
}

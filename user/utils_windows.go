package user

import (
	"sort"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

func getGroupAttrs(attributes uint32) ([]string, error) {
	var attrs []string
	var keys []int
	var mask map[int]string = map[int]string{
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
	var valid int = 0xe000007f

	if int(attributes)|valid != valid {
		return nil, nil
	}

	for k := range mask {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		if int(attributes)&k > 0 {
			attrs = append(attrs, mask[k])
		}
	}

	return attrs, nil
}

func getGroupNameAndType(sid *windows.SID) (string, string, error) {
	var account string
	var acctype string
	var domain string
	var name string
	var theType uint32

	account, domain, theType, _ = sid.LookupAccount(".")

	if account == "None" {
		return "", "", nil
	}

	if strings.HasPrefix(account, "LogonSessionId_") {
		return "", "", nil
	}

	name = domain + "\\" + account
	if domain == "" {
		name = account
	}

	switch theType {
	case 0:
		acctype = "Unknown SID type"
	case 1:
		acctype = "user"
	case 2:
		acctype = "Group"
	case 3:
		acctype = "Domain"
	case 4:
		acctype = "Alias"
	case 5:
		acctype = "Well-known group"
	case 6:
		acctype = "Deleted"
	case 7:
		acctype = "Invalid"
	case 8:
		acctype = "Computer"
	case 10:
		acctype = "Label"
	}

	return name, acctype, nil
}

func getPrivName(luid uint64) (string, error) {
	var b []byte
	var e error
	var n int

	if e = w32.LookupPrivilegeName("", luid, &b, &n); e != nil {
		b = make([]byte, n)
		if e = w32.LookupPrivilegeName("", luid, &b, &n); e != nil {
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
	var line string
	var lines []string
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
	line = ""
	for i, col := range hdrs {
		line += col + strings.Repeat(" ", width[i]-len(col)) + " "
	}

	lines = append(lines, line)

	// Dividers
	line = ""
	for i := range hdrs {
		line += strings.Repeat("=", width[i]) + " "
	}

	lines = append(lines, line)

	// Data
	for _, row := range data {
		line = ""

		for i, col := range row {
			line += col + strings.Repeat(" ", width[i]-len(col)) + " "
		}

		lines = append(lines, line)
	}

	// Print
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}

func tokenOrDefault(access []windows.Token) windows.Token {
	if len(access) == 0 {
		return windows.GetCurrentProcessToken()
	}

	return access[0]
}

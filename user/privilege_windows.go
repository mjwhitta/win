package user

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

// Privilege contains information about a Windows privilege.
type Privilege struct {
	Description string
	Name        string
	State       string
}

// Privileges returns an array of Privileges associated with the
// provided access token. If no token is provided, it defaults to the
// current process.
func Privileges(access ...windows.Token) ([]Privilege, error) {
	var attrs uint32
	var b []byte
	var buf *bytes.Buffer
	var e error
	var n uint32
	var luid uint64
	var privs []Privilege

	// Get number of bytes
	windows.GetTokenInformation(
		tokenOrDefault(access),
		windows.TokenPrivileges,
		nil,
		0,
		&n,
	)

	// Now create memory and fill it in
	b = make([]byte, n)
	e = windows.GetTokenInformation(
		tokenOrDefault(access),
		windows.TokenPrivileges,
		&b[0],
		n,
		&n,
	)
	if e != nil {
		e = errors.Newf("failed to get token privileges: %w", e)
		return nil, e
	}

	// Read number of privileges
	buf = bytes.NewBuffer(b)
	if e = binary.Read(buf, binary.LittleEndian, &n); e != nil {
		e = errors.Newf("failed to read number of privileges: %w", e)
		return nil, e
	}

	privs = make([]Privilege, n)

	for i := range privs {
		e = binary.Read(buf, binary.LittleEndian, &luid)
		if e != nil {
			return nil, errors.Newf("failed to read LUID: %w", e)
		}

		if privs[i].Name, e = getPrivName(luid); e != nil {
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

		privs[i].State = "Disabled"
		if (attrs & windows.SE_PRIVILEGE_ENABLED) > 0 {
			privs[i].State = "Enabled"
		}
	}

	return privs, nil
}

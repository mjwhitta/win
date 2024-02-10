package user

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

// Privilege contains information about a Windows privilege.
type Privilege struct {
	Attributes  uint32
	Description string
	LUID        uint64
	Name        string
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

		privs[i].LUID = luid

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

		privs[i].Attributes = attrs
	}

	return privs, nil
}

// Enabled will return whether or not the Privilege has been enabled.
func (p *Privilege) Enabled() bool {
	return p.Attributes&windows.SE_PRIVILEGE_ENABLED > 0
}

// EnabledByDefault will return whether or not the Privilege is
// enabled by default.
func (p *Privilege) EnabledByDefault() bool {
	return p.Attributes&windows.SE_PRIVILEGE_ENABLED_BY_DEFAULT > 0
}

// Removed will return whether or not the Privilege has been removed.
func (p *Privilege) Removed() bool {
	return p.Attributes&windows.SE_PRIVILEGE_REMOVED > 0
}

// UsedForAccess will return whether or not the Privilege is used for
// access.
func (p *Privilege) UsedForAccess() bool {
	return p.Attributes&windows.SE_PRIVILEGE_USED_FOR_ACCESS > 0
}

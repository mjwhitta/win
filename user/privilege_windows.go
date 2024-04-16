package user

import (
	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

// Privilege contains information about a Windows privilege.
type Privilege struct {
	Attributes  uint32
	Description string
	LUID        windows.LUID
	Name        string

	proc windows.Handle
}

// Privileges returns an array of Privileges for the process token
// associated with the provided process handle. If no handle is
// provided, it defaults to the current process.
func Privileges(proc ...windows.Handle) ([]*Privilege, error) {
	var b []byte
	var e error
	var n uint32
	var t windows.Token = tokenOrDefault(proc)

	// Get number of bytes
	windows.GetTokenInformation(
		t,
		windows.TokenPrivileges,
		nil,
		0,
		&n,
	)

	// Now create memory and fill it in
	b = make([]byte, n)
	e = windows.GetTokenInformation(
		t,
		windows.TokenPrivileges,
		&b[0],
		n,
		&n,
	)
	if e != nil {
		e = errors.Newf("failed to get token privileges: %w", e)
		return nil, e
	}

	// Read privileges from bytes
	return privsFromBytes(b, n, procOrDefault(proc))
}

// Disable will adjust token privileges to disable the Privilege.
func (p *Privilege) Disable() error {
	if !p.Enabled() {
		return nil
	}

	p.Attributes ^= windows.SE_PRIVILEGE_ENABLED
	return adjustToken(p)
}

// Enable will adjust token privileges to enable the Privilege.
func (p *Privilege) Enable() error {
	if p.Enabled() {
		return nil
	}

	p.Attributes ^= windows.SE_PRIVILEGE_ENABLED
	return adjustToken(p)
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

// Remove will adjust token privileges to remove the Privilege.
func (p *Privilege) Remove() error {
	if p.Removed() {
		return nil
	}

	p.Attributes ^= windows.SE_PRIVILEGE_REMOVED
	return adjustToken(p)
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

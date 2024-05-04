//go:build windows

package user

import (
	"strings"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

// Name will return the name of the current user.
func Name() (string, error) {
	var e error
	var name *uint16
	var size uint32
	var tmp []uint16

	windows.GetUserNameEx(windows.NameSamCompatible, name, &size)

	if size > 0 {
		tmp = make([]uint16, size)
		name = &tmp[0]
	}

	e = windows.GetUserNameEx(windows.NameSamCompatible, name, &size)
	if e != nil {
		return "", errors.Newf("failed to get username: %w", e)
	}

	return strings.ToLower(windows.UTF16PtrToString(name)), nil
}

// SID will return the SID of the user for the process token
// associated with the provided process handle. If no handle is
// provided, it defaults to the current process.
func SID(proc ...windows.Handle) (string, error) {
	var e error
	var tu *windows.Tokenuser

	if tu, e = tokenOrDefault(proc).GetTokenUser(); e != nil {
		return "", errors.Newf("failed to get token user: %w", e)
	}

	return tu.User.Sid.String(), nil
}

// Whoami will return output that very nearly (if not exactly) matches
// the "whoami.exe" output. If no process handle is provided, it
// defaults to the current process.
func Whoami(proc ...windows.Handle) (string, error) {
	var e error
	var id *ID

	if id, e = Identity(proc...); e != nil {
		return "", e
	}

	return id.Whoami(), nil
}

// WhoamiAll will return output that very nearly (if not exactly)
// matches the "whoami.exe /all" output. If no process handle is
// provided, it defaults to the current process.
func WhoamiAll(proc ...windows.Handle) (string, error) {
	var e error
	var id *ID

	if id, e = Identity(proc...); e != nil {
		return "", e
	}

	return id.WhoamiAll(), nil
}

// WhoamiGroups will return output that very nearly (if not exactly)
// matches the "whoami.exe /groups" output. If no process handle is
// provided, it defaults to the current process.
func WhoamiGroups(proc ...windows.Handle) (string, error) {
	var e error
	var id *ID

	if id, e = Identity(proc...); e != nil {
		return "", e
	}

	return id.WhoamiGroups(), nil
}

// WhoamiPriv will return output that very nearly (if not exactly)
// matches the "whoami.exe /priv" output. If no process handle is
// provided, it defaults to the current process.
func WhoamiPriv(proc ...windows.Handle) (string, error) {
	var e error
	var id *ID

	if id, e = Identity(proc...); e != nil {
		return "", e
	}

	return id.WhoamiPriv(), nil
}

// WhoamiUser will return output that very nearly (if not exactly)
// matches the "whoami.exe /user" output. If no process handle is
// provided, it defaults to the current process.
func WhoamiUser(proc ...windows.Handle) (string, error) {
	var e error
	var id *ID

	if id, e = Identity(proc...); e != nil {
		return "", e
	}

	return id.WhoamiUser(), nil
}

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

// SID will return the SID for the user associated with the provided
// access token. If no token is provided, it defaults to the current
// process.
func SID(access ...windows.Token) (string, error) {
	var e error
	var tu *windows.Tokenuser

	if tu, e = tokenOrDefault(access).GetTokenUser(); e != nil {
		return "", errors.Newf("failed to get token user: %w", e)
	}

	return tu.User.Sid.String(), nil
}

// Whoami will return output that very nearly (if not exactly) matches
// the "whoami.exe /all" output. If no access token is provided, it
// defaults to the current process.
func Whoami(access ...windows.Token) (string, error) {
	var e error
	var groups [][]string
	var id *ID
	var lines []string
	var privs [][]string
	var state string

	if id, e = Identity(tokenOrDefault(access)); e != nil {
		return "", e
	}

	// User info
	lines = append(
		lines,
		output(
			"USER INFORMATION",
			[]string{"User Name", "SID"},
			[][]string{{id.Name, id.SID}},
		),
	)

	// Group info
	for _, group := range id.Groups {
		groups = append(
			groups,
			[]string{
				group.Name,
				group.Type,
				group.SID,
				strings.Join(group.Attrs, ", "),
			},
		)
	}

	lines = append(
		lines,
		output(
			"GROUP INFORMATION",
			[]string{"Group Name", "Type", "SID", "Attributes"},
			groups,
		),
	)

	// Privileges info
	for _, priv := range id.Privileges {
		state = "Disabled"
		if priv.Enabled() {
			state = "Enabled"
		}

		privs = append(
			privs,
			[]string{priv.Name, priv.Description, state},
		)
	}

	lines = append(
		lines,
		output(
			"PRIVILEGES INFORMATION",
			[]string{"Privilege Name", "Description", "State"},
			privs,
		),
	)

	return strings.Join(lines, "\n"), nil
}

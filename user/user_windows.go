package user

import (
	"strings"

	"github.com/mjwhitta/errors"
	"golang.org/x/sys/windows"
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
// process token. If no token is provided, it defaults to the current
// process.
func SID(processToken ...windows.Token) (string, error) {
	var e error
	var t windows.Token
	var tu *windows.Tokenuser

	if len(processToken) == 0 {
		t = windows.GetCurrentProcessToken()
	} else {
		t = processToken[0]
	}

	if tu, e = t.GetTokenUser(); e != nil {
		return "", errors.Newf("failed to get token user: %w", e)
	}

	return tu.User.Sid.String(), nil
}

// Whoami will return output that very nearly (if not exactly) matches
// the "whoami.exe /all" output.
func Whoami() (string, error) {
	var e error
	var groups [][]string
	var id *ID
	var lines []string

	if id, e = Identity(); e != nil {
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

	return strings.Join(lines, "\n"), nil
}

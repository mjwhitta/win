package user

import (
	"strings"

	"golang.org/x/sys/windows"
)

// ID contains information about a Windows user.
type ID struct {
	Groups     []Group
	Name       string
	Privileges []Privilege
	SID        string
}

// Identity will return a pointer to a new ID instance containing the
// user information for the provided access token. If no token is
// provided, it defaults to the current process.
func Identity(access ...windows.Token) (*ID, error) {
	var e error
	var groups []Group
	var id *ID
	var name string
	var privs []Privilege
	var sid string

	if name, e = Name(); e != nil {
		return nil, e
	}

	if sid, e = SID(tokenOrDefault(access)); e != nil {
		return nil, e
	}

	if groups, e = Groups(tokenOrDefault(access)); e != nil {
		return nil, e
	}

	if privs, e = Privileges(tokenOrDefault(access)); e != nil {
		return nil, e
	}

	id = &ID{Groups: groups, Name: name, Privileges: privs, SID: sid}
	return id, nil
}

// HasPrivilege will search the associated Privileges for one with the
// provided name. The returned bool should be checked before using the
// returned Privilege.
func (id *ID) HasPrivilege(name string) (Privilege, bool) {
	for _, priv := range id.Privileges {
		if priv.Name == name {
			return priv, true
		}
	}

	return Privilege{}, false
}

// InGroup will search the associated Groups for one with the provided
// name. The returned bool should be checked before using the returned
// Group.
func (id *ID) InGroup(name string) (Group, bool) {
	for _, group := range id.Groups {
		if group.Name == name {
			return group, true
		}
	}

	return Group{}, false
}

// Whoami will return output that very nearly (if not exactly) matches
// the "whoami.exe" output.
func (id *ID) Whoami() string {
	return id.Name
}

// WhoamiAll will return output that very nearly (if not exactly)
// matches the "whoami.exe /all" output.
func (id *ID) WhoamiAll() string {
	return strings.Join(
		[]string{
			id.WhoamiUser(),
			"",
			id.WhoamiGroups(),
			"",
			id.WhoamiPriv(),
			"",
		},
		"\n",
	)
}

// WhoamiGroups will return output that very nearly (if not exactly)
// matches the "whoami.exe /groups" output.
func (id *ID) WhoamiGroups() string {
	var groups [][]string
	var lines []string

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

	return strings.Join(lines, "\n")
}

// WhoamiPriv will return output that very nearly (if not exactly)
// matches the "whoami.exe /priv" output.
func (id *ID) WhoamiPriv() string {
	var lines []string
	var privs [][]string
	var state string

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

	return strings.Join(lines, "\n")
}

// WhoamiUser will return output that very nearly (if not exactly)
// matches the "whoami.exe /user" output.
func (id *ID) WhoamiUser() string {
	var lines []string

	// User info
	lines = append(
		lines,
		output(
			"USER INFORMATION",
			[]string{"User Name", "SID"},
			[][]string{{id.Name, id.SID}},
		),
	)

	return strings.Join(lines, "\n")
}

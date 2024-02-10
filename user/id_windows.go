package user

import "golang.org/x/sys/windows"

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
func (i *ID) HasPrivilege(name string) (Privilege, bool) {
	for _, priv := range i.Privileges {
		if priv.Name == name {
			return priv, true
		}
	}

	return Privilege{}, false
}

// InGroup will search the associated Groups for one with the provided
// name. The returned bool should be checked before using the returned
// Group.
func (i *ID) InGroup(name string) (Group, bool) {
	for _, group := range i.Groups {
		if group.Name == name {
			return group, true
		}
	}

	return Group{}, false
}

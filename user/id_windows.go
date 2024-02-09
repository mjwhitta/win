package user

import "golang.org/x/sys/windows"

// ID contains information about a Windows user.
type ID struct {
	Groups []Group
	Name   string
	SID    string
}

// Identity will return a pointer to a new ID instance containing the
// user information for the provided process token. If no token is
// provided, it defaults to the current process.
func Identity(processToken ...windows.Token) (*ID, error) {
	var e error
	var groups []Group
	var name string
	var sid string
	var t windows.Token

	if len(processToken) == 0 {
		t = windows.GetCurrentProcessToken()
	} else {
		t = processToken[0]
	}

	if name, e = Name(); e != nil {
		return nil, e
	}

	if sid, e = SID(t); e != nil {
		return nil, e
	}

	if groups, e = Groups(t); e != nil {
		return nil, e
	}

	return &ID{Groups: groups, Name: name, SID: sid}, nil
}

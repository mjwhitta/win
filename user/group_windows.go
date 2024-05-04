//go:build windows

package user

import (
	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

// Group contains information about a Windows group.
type Group struct {
	Attrs []string
	Name  string
	SID   string
	Type  string
}

// Groups returns an array of Groups for the process token associated
// with the provided process handle. If no handle is provided, it
// defaults to the current process.
func Groups(proc ...windows.Handle) ([]Group, error) {
	var acctype string
	var attrs []string
	var e error
	var groups []Group
	var name string
	var tg *windows.Tokengroups

	if tg, e = tokenOrDefault(proc).GetTokenGroups(); e != nil {
		e = errors.Newf("failed to get token groups: %w", e)
		return nil, e
	}

	for _, g := range tg.AllGroups() {
		if name, acctype, e = getGroupNameAndType(g.Sid); e != nil {
			return nil, e
		}

		if acctype == "" {
			continue
		}

		attrs = []string{}
		if acctype != "Label" {
			if attrs, e = getGroupAttrs(g.Attributes); e != nil {
				return nil, e
			}
		}

		groups = append(
			groups,
			Group{
				Attrs: attrs,
				Name:  name,
				SID:   g.Sid.String(),
				Type:  acctype,
			},
		)
	}

	return groups, nil
}

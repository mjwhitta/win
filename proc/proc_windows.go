//go:build windows

package proc

import (
	"sort"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

// Get will return a list of Windows processes.
func Get() ([]*w32.ProcessEntry32, error) {
	var e error
	var pe *w32.ProcessEntry32
	var procs []*w32.ProcessEntry32
	var snapHndl windows.Handle

	snapHndl, e = w32.CreateToolhelp32Snapshot(
		w32.Tlhelp32.SnapProcess,
		0,
	)
	if e != nil {
		return nil, errors.Newf("failed to create snapshot: %w", e)
	}
	defer windows.Close(snapHndl)

	if pe, e = w32.Process32First(snapHndl); e != nil {
		return nil, errors.Newf("failed to get first process: %w", e)
	}

	for {
		if pe == nil {
			break
		}

		procs = append(procs, pe)

		if pe, e = w32.Process32Next(snapHndl); e != nil {
			e = errors.Newf("failed to get next process: %w", e)
			return nil, e
		}
	}

	sort.Slice(procs, psLessFunc(procs))
	return procs, nil
}

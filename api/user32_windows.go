//go:build windows

package api

import (
	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var user32 *windows.LazyDLL = windows.NewLazySystemDLL("user32")

// CreateDesktopW from winuser.h
func CreateDesktopW(
	name string, flags uintptr, access uintptr,
) (windows.Handle, error) {
	var desktop uintptr
	var e error
	var proc string = "CreateDesktopW"

	desktop, _, e = user32.NewProc(proc).Call(
		types.LpCwstr(name),
		0,
		0,
		flags,
		access,
		0,
	)
	if desktop == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return windows.Handle(desktop), nil
}

// SwitchDesktop from winuser.h
func SwitchDesktop(desktop windows.Handle) error {
	var e error
	var proc string = "SwitchDesktop"
	var success uintptr

	success, _, e = user32.NewProc(proc).Call(uintptr(desktop))
	if success == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

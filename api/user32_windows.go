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
	name string,
	flags uintptr,
	access uintptr,
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

// EnumWindowStationsW from winuser.h
//
// This function accepts an enumeration function, which must be
// declared in C with cgo, but can optionally be defined in Go.
//
// To declare your callback function in C with cgo:
//
//	// unsigned char cb(short unsigned int*, long long unsigned int);
//	import "C"
//
// To optionally define in Go (otherwise define in C above):
//
//	//export cb
//	func cb(name *uint16, params uintptr) bool {
//		println(windows.UTF16PtrToString(name))
//		return true
//	}
//
// To call this function in Go:
//
//	func main() {
//		_ = w32.EnumWindowStationsW(uintptr(unsafe.Pointer(C.cb)), 0)
//	}
func EnumWindowStationsW(enumFunc uintptr, params uintptr) error {
	var e error
	var proc string = "EnumWindowStationsW"
	var success uintptr

	success, _, _ = user32.NewProc(proc).Call(enumFunc, params)
	if success == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
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

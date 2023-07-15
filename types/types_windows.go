package types

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Cwstr converts a Go string to a Windows wide string.
func Cwstr(str string) *uint16 {
	var tmp *uint16

	tmp, _ = windows.UTF16PtrFromString(str)
	return tmp
}

// LpCwstr converts a Go string to a Windows wide string pointer.
func LpCwstr(str string) uintptr {
	if str == "" {
		return 0
	}

	return uintptr(unsafe.Pointer(Cwstr(str)))
}

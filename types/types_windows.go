//go:build windows

package types

import (
	"unsafe"

	"github.com/mjwhitta/errors"
	"golang.org/x/sys/windows"
)

// Cwstr converts a Go string to a Windows wide string.
func Cwstr(str string) (tmp *uint16) {
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

// Utf16LE will convert a Go string to a utf16 array that is little
// endian byte-ordered. It will then return a Go []byte. The output
// can then be base64 encoded for use with PowerShell's encoded
// command functionality.
func Utf16LE(str string) ([]byte, error) {
	var b []byte
	var e error
	var u []uint16

	if u, e = windows.UTF16FromString(str); e != nil {
		return nil, errors.Newf("failed to create wide string: %w", e)
	}

	for i, c := range u {
		if i == len(u)-1 {
			break
		}

		b = append(b, byte(c), '\x00')
	}

	return b, nil
}

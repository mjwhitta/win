//go:build windows

package api

import (
	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var kernel32 *windows.LazyDLL = windows.NewLazySystemDLL("kernel32")

// HeapAlloc is from heapapi.h
func HeapAlloc(
	heapHndl uintptr,
	dwFlags uintptr,
	dwBytes int,
) (uintptr, error) {
	var e error
	var addr uintptr
	var proc string = "HeapAlloc"

	addr, _, e = kernel32.NewProc(proc).Call(
		heapHndl,
		dwFlags,
		uintptr(dwBytes),
	)
	if addr == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return addr, nil
}

// HeapCreate is from heapapi.h
func HeapCreate(
	flOptions uintptr,
	dwInitialSize int,
	dwMaximumSize int,
) (uintptr, error) {
	var e error
	var hndl uintptr
	var proc string = "HeapCreate"

	hndl, _, e = kernel32.NewProc(proc).Call(
		flOptions,
		uintptr(dwInitialSize),
		uintptr(dwMaximumSize),
	)
	if hndl == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return hndl, nil
}

// HeapDestroy is from heapapi.h
func HeapDestroy(hndl uintptr) error {
	var e error
	var ok uintptr
	var proc string = "HeapDestroy"

	if ok, _, e = kernel32.NewProc(proc).Call(hndl); ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// HeapFree is from heapapi.h
func HeapFree(heapHndl uintptr, dwFlags uintptr, addr uintptr) error {
	var e error
	var ok uintptr
	var proc string = "HeapFree"

	ok, _, e = kernel32.NewProc(proc).Call(heapHndl, dwFlags, addr)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// OutputDebugStringW will print a string that Dbgview.exe and
// dbgview64.exe will display. Useful for debugging DLLs.
func OutputDebugStringW(out string) {
	var proc string = "OutputDebugStringW"

	kernel32.NewProc(proc).Call(types.LpCwstr(out))
}

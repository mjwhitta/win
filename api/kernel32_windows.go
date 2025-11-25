//go:build windows

package api

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var kernel32 *windows.LazyDLL = windows.NewLazySystemDLL("kernel32")

// CopyFile2 from winbase.h
func CopyFile2(
	src string,
	dst string,
	params CopyFile2ExtendedParameters,
) error {
	var e error
	var ok uintptr
	var proc string = "CopyFile2"

	params.dwSize = uint32(unsafe.Sizeof(params)) // Always 32

	ok, _, e = kernel32.NewProc(proc).Call(
		types.LpCwstr(src),
		types.LpCwstr(dst),
		uintptr(unsafe.Pointer(&params)),
	)
	if ok != 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// CreateToolhelp32Snapshot from tlhelp32.h
func CreateToolhelp32Snapshot(
	dwFlags uintptr,
	th32ProcessID uintptr,
) (windows.Handle, error) {
	var e error
	var hndl uintptr
	var proc string = "CreateToolhelp32Snapshot"

	hndl, _, e = kernel32.NewProc(proc).Call(dwFlags, th32ProcessID)
	if hndl == 0 {
		return 0, errors.Newf("%s: %s", proc, e.Error())
	}

	return windows.Handle(hndl), nil
}

// HeapAlloc from heapapi.h
func HeapAlloc(
	heapHndl uintptr,
	dwFlags uintptr,
	dwBytes uintptr,
) (uintptr, error) {
	var e error
	var addr uintptr
	var proc string = "HeapAlloc"

	addr, _, e = kernel32.NewProc(proc).Call(
		heapHndl,
		dwFlags,
		dwBytes,
	)
	if addr == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return addr, nil
}

// HeapCreate from heapapi.h
func HeapCreate(
	flOptions uintptr,
	dwInitialSize uintptr,
	dwMaximumSize uintptr,
) (uintptr, error) {
	var e error
	var hndl uintptr
	var proc string = "HeapCreate"

	hndl, _, e = kernel32.NewProc(proc).Call(
		flOptions,
		dwInitialSize,
		dwMaximumSize,
	)
	if hndl == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return hndl, nil
}

// HeapDestroy from heapapi.h
func HeapDestroy(hndl uintptr) error {
	var e error
	var ok uintptr
	var proc string = "HeapDestroy"

	if ok, _, e = kernel32.NewProc(proc).Call(hndl); ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// HeapFree from heapapi.h
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

	_, _, _ = kernel32.NewProc(proc).Call(types.LpCwstr(out))
}

// Process32First from tlhelp32.h
func Process32First(
	snapHndl windows.Handle,
) (*ProcessEntry32, error) {
	var e error
	var ok uintptr
	var pe ProcessEntry32
	var proc string = "Process32First"

	pe.dwSize = uint32(unsafe.Sizeof(pe)) // Always 304

	ok, _, e = kernel32.NewProc(proc).Call(
		uintptr(snapHndl),
		uintptr(unsafe.Pointer(&pe)),
	)
	if ok == 0 {
		if strings.Contains(e.Error(), "There are no more files") {
			//nolint:nilnil // Not a real error, but we are done
			return nil, nil
		}

		return nil, errors.Newf("%s: %s", proc, e.Error())
	}

	return &pe, nil
}

// Process32Next from tlhelp32.h
func Process32Next(snapHndl windows.Handle) (*ProcessEntry32, error) {
	var e error
	var ok uintptr
	var pe ProcessEntry32
	var proc string = "Process32Next"

	pe.dwSize = uint32(unsafe.Sizeof(pe)) // Always 304

	ok, _, e = kernel32.NewProc(proc).Call(
		uintptr(snapHndl),
		uintptr(unsafe.Pointer(&pe)),
	)
	if ok == 0 {
		if strings.Contains(e.Error(), "There are no more files") {
			//nolint:nilnil // Not a real error, but we are done
			return nil, nil
		}

		return nil, errors.Newf("%s: %s", proc, e.Error())
	}

	return &pe, nil
}

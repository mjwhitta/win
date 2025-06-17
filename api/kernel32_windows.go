//go:build windows

package api

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

// CopyFile2ExtendedParameters is COPYFILE2_EXTENDED_PARAMETERS from
// winbase.h
type CopyFile2ExtendedParameters struct {
	dwSize          uint32  // DWORD, 4 bytes, always 32
	CopyFlags       uint32  // DWORD, 4 bytes
	CancelPtr       uintptr // pointer, 8 bytes
	ProgressRoutine uintptr // pointer, 8 bytes
	CallbackContext uintptr // pointer, 8 bytes
}

// ProcessEntry32 is PROCESSENTRY32 from tlhelp32.h
type ProcessEntry32 struct {
	dwSize            uint32     // DWORD, 4 bytes, always 304
	cntUsage          uint32     // DWORD, 4 bytes, always 0
	PID               uint32     // DWORD, 4 bytes (+4 for alignment?)
	defaultHeapID     uintptr    // pointer, 8 bytes, always 0
	moduleID          uint32     // DWORD, 4 bytes, always 0
	ThreadCount       uint32     // DWORD, 4 bytes
	ParentPID         uint32     // DWORD, 4 bytes
	PriorityClassBase uint32     // LONG, 4 bytes
	dwFlags           uint32     // DWORD, 4 bytes, always 0
	exeFile           [260]uint8 // char*, Stdlib.MaxPath bytes
}

// ExeFile will convert the exe filename to a Go string.
func (pe *ProcessEntry32) ExeFile() string {
	return windows.ByteSliceToString(pe.exeFile[:])
}

var kernel32 *windows.LazyDLL = windows.NewLazySystemDLL("kernel32")

// CopyFile2 from winbase.h
func CopyFile2(
	src string, dst string, params CopyFile2ExtendedParameters,
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

// HeapCreate from heapapi.h
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

	kernel32.NewProc(proc).Call(types.LpCwstr(out))
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
			return nil, nil
		}

		return nil, errors.Newf("%s: %s", proc, e.Error())
	}

	return &pe, nil
}

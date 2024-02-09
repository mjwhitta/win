package api

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var advapi32 *windows.LazyDLL = windows.NewLazySystemDLL("Advapi32")

// LookupPrivilegeDisplayName is from winbase.h
func LookupPrivilegeDisplayName(
	system string, name string, desc *[]byte, descLen *int,
) error {
	var b []uint16
	var e error
	var langID uint32
	var proc string = "LookupPrivilegeDisplayNameW"
	var success uintptr
	var tmp string

	if *descLen > 0 {
		b = make([]uint16, *descLen)
	} else {
		b = make([]uint16, 1)
	}

	success, _, e = advapi32.NewProc(proc).Call(
		types.LpCwstr(system),
		types.LpCwstr(name),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(descLen)),
		uintptr(unsafe.Pointer(&langID)),
	)
	if success == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	tmp = windows.UTF16ToString(b)
	*desc = []byte(tmp)

	return nil
}

// LookupPrivilegeName is from winbase.h
func LookupPrivilegeName(
	system string, luid uint64, name *[]byte, nameLen *int,
) error {
	var b []uint16
	var e error
	var proc string = "LookupPrivilegeNameW"
	var success uintptr
	var tmp string

	if *nameLen > 0 {
		b = make([]uint16, *nameLen)
	} else {
		b = make([]uint16, 1)
	}

	success, _, e = advapi32.NewProc(proc).Call(
		types.LpCwstr(system),
		uintptr(unsafe.Pointer(&luid)),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(nameLen)),
	)
	if success == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	tmp = windows.UTF16ToString(b)
	*name = []byte(tmp)

	return nil
}

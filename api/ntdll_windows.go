//go:build windows

package api

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
)

type clientID struct {
	UniqueProcess uintptr
	UniqueThread  uintptr
}

type objectAttrs struct {
	Length                   uintptr
	RootDirectory            uintptr
	ObjectName               uintptr
	Attributes               uintptr
	SecurityDescriptor       uintptr
	SecurityQualityOfService uintptr
}

var ntdll *windows.LazyDLL = windows.NewLazySystemDLL("ntdll")

// NtAllocateVirtualMemory from ntdll.
func NtAllocateVirtualMemory(
	pHndl windows.Handle,
	size uint64,
	allocType uintptr,
	protection uintptr,
) (uintptr, error) {
	var addr uintptr
	var err uintptr
	var proc string = "NtAllocateVirtualMemory"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(pHndl),
		uintptr(unsafe.Pointer(&addr)),
		0,
		uintptr(unsafe.Pointer(&size)),
		allocType,
		protection,
	)
	if err != 0 {
		return 0, errors.Newf("%s returned %0x", proc, uint32(err))
	} else if addr == 0 {
		return 0, errors.Newf("%s failed for unknown reason", proc)
	}

	// WTF?! Why is a Printf needed?! time.Sleep() doesn't work?
	// Printf("") doesn't work? Oh well, print newline and escape
	// sequence for "go up 1 line"
	fmt.Printf("\n\x1b[1A")

	return addr, nil
}

// NtCreateSection from ntdll.
func NtCreateSection(
	sHndl *windows.Handle,
	access uintptr,
	size uint64,
	pagePerms uintptr,
	secPerms uintptr,
) error {
	var err uintptr
	var proc string = "NtCreateSection"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(unsafe.Pointer(sHndl)),
		access,
		0,
		uintptr(unsafe.Pointer(&size)),
		pagePerms,
		secPerms,
		0,
	)
	if err != 0 {
		return errors.Newf("%s returned %0x", proc, uint32(err))
	} else if *sHndl == 0 {
		return errors.Newf("%s failed for unknown reason", proc)
	}

	return nil
}

// NtCreateThreadEx from ntdll.
func NtCreateThreadEx(
	pHndl windows.Handle,
	addr uintptr,
	sspnd bool,
) (windows.Handle, error) {
	var err uintptr
	var proc string = "NtCreateThreadEx"
	var suspend uintptr
	var tHndl windows.Handle

	if sspnd {
		suspend = 1
	}

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(unsafe.Pointer(&tHndl)),
		Winnt.ThreadAllAccess,
		0,
		uintptr(pHndl),
		addr,
		0,
		suspend,
		0,
		0,
		0,
		0,
	)
	if err != 0 {
		return 0, errors.Newf("%s returned %0x", proc, uint32(err))
	} else if tHndl == 0 {
		return 0, errors.Newf("%s failed for unknown reason", proc)
	}

	return tHndl, nil
}

// NtMapViewOfSection from ntdll.
func NtMapViewOfSection(
	sHndl windows.Handle,
	pHndl windows.Handle,
	size uint64,
	inheritPerms uintptr,
	pagePerms uintptr,
) (uintptr, error) {
	var err uintptr
	var proc string = "NtMapViewOfSection"
	var scBase uintptr
	var scOffset uintptr

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(sHndl),
		uintptr(pHndl),
		uintptr(unsafe.Pointer(&scBase)),
		0,
		0,
		uintptr(unsafe.Pointer(&scOffset)),
		uintptr(unsafe.Pointer(&size)),
		inheritPerms,
		0,
		pagePerms,
	)
	if err != 0 {
		return 0, errors.Newf("%s returned %0x", proc, uint32(err))
	} else if scBase == 0 {
		return 0, errors.Newf("%s failed for unknown reason", proc)
	}

	return scBase, nil
}

// NtOpenProcess from ntdll.
func NtOpenProcess(
	pid uint32,
	access uintptr,
) (windows.Handle, error) {
	var err uintptr
	var pHndl windows.Handle
	var proc string = "NtOpenProcess"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(unsafe.Pointer(&pHndl)),
		access,
		uintptr(unsafe.Pointer(&objectAttrs{0, 0, 0, 0, 0, 0})),
		uintptr(unsafe.Pointer(&clientID{uintptr(pid), 0})),
	)
	if err != 0 {
		return 0, errors.Newf("%s returned %0x", proc, uint32(err))
	} else if pHndl == 0 {
		return 0, errors.Newf("%s failed for unknown reason", proc)
	}

	return pHndl, nil
}

// NtQueueApcThread from ntdll.
func NtQueueApcThread(
	tHndl windows.Handle,
	apcRoutine uintptr,
) error {
	var err uintptr
	var proc string = "NtQueueApcThread"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(tHndl),
		apcRoutine,
		0, // arg1
		0, // arg2
		0, // arg3
	)
	if err != 0 {
		return errors.Newf("%s returned %0x", proc, uint32(err))
	}

	return nil
}

// NtQueueApcThreadEx from ntdll.
func NtQueueApcThreadEx(
	tHndl windows.Handle,
	apcRoutine uintptr,
) error {
	var err uintptr
	var proc string = "NtQueueApcThreadEx"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(tHndl),
		0x1, // userApcReservedHandle
		apcRoutine,
		0, // arg1
		0, // arg2
		0, // arg3
	)
	if err != 0 {
		return errors.Newf("%s returned %0x", proc, uint32(err))
	}

	return nil
}

// NtResumeThread from ntdll.
func NtResumeThread(tHndl windows.Handle) error {
	var err uintptr
	var proc string = "NtResumeThread"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(tHndl),
		0, // previousSuspendCount
	)
	if err != 0 {
		return errors.Newf("%s returned %0x", proc, uint32(err))
	}

	return nil
}

// NtWriteVirtualMemory from ntdll.
func NtWriteVirtualMemory(
	pHndl windows.Handle,
	dst uintptr,
	b []byte,
) error {
	var err uintptr
	var proc string = "NtWriteVirtualMemory"

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(pHndl),
		dst,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
	)
	if err != 0 {
		return errors.Newf("%s returned %0x", proc, uint32(err))
	}

	return nil
}

// RtlCreateUserThread from ntdll.
func RtlCreateUserThread(
	pHndl windows.Handle,
	addr uintptr,
	sspnd bool,
) (windows.Handle, error) {
	var err uintptr
	var proc string = "RtlCreateUserThread"
	var suspend uintptr
	var tHndl windows.Handle

	if sspnd {
		suspend = 1
	}

	err, _, _ = ntdll.NewProc(proc).Call(
		uintptr(pHndl),
		0,
		suspend,
		0,
		0,
		0,
		addr,
		0,
		uintptr(unsafe.Pointer(&tHndl)),
		0,
	)
	if err != 0 {
		return 0, errors.Newf("%s returned %0x", proc, uint32(err))
	} else if tHndl == 0 {
		return 0, errors.Newf("%s failed for unknown reason", proc)
	}

	return tHndl, nil
}

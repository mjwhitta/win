//go:build windows

package api

// CopyFile2ExtendedParameters is COPYFILE2_EXTENDED_PARAMETERS from
// winbase.h
type CopyFile2ExtendedParameters struct {
	dwSize          uint32  // DWORD, 4 bytes, always 32
	CopyFlags       uint32  // DWORD, 4 bytes
	CancelPtr       uintptr // pointer, 8 bytes
	ProgressRoutine uintptr // pointer, 8 bytes
	CallbackContext uintptr // pointer, 8 bytes
}

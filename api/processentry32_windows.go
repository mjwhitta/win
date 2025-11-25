//go:build windows

package api

import "golang.org/x/sys/windows"

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

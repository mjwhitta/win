//go:build windows

package main

// unsigned char cb(short unsigned int*, long long unsigned int);
import "C"

import (
	"unsafe"

	w32 "github.com/mjwhitta/win/api"
	"golang.org/x/sys/windows"
)

//export cb
func cb(name *uint16, _ uintptr) bool {
	println(windows.UTF16PtrToString(name))
	return true
}

func main() {
	_ = w32.EnumWindowStationsW(uintptr(unsafe.Pointer(C.cb)), 0)
}

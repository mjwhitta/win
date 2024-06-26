//go:build windows

package api

// unsigned char enumWindowStationsHook(
//   short unsigned int* lpszWindowStation,
//   long long unsigned int lParam
// );
import "C"

import (
	"sync"
	"unsafe"

	"github.com/mjwhitta/errors"
	"golang.org/x/sys/windows"
)

var enumWS struct {
	fp func(string, uintptr) bool
	m  sync.Mutex
}

// EnumWindowStationsW from winuser.h
//
// This function accepts an enumeration function, which can be either
// of the following:
//
// - func(string, uintptr) bool <- a Go func
// - uintptr <- a poter to C function address
//
// Anything else returns an errors.
func EnumWindowStationsW(enumFunc any, params uintptr) error {
	var proc string = "EnumWindowStationsW"

	switch enumFunc := enumFunc.(type) {
	case func(string, uintptr) bool:
		// This function isn't thread-safe b/c we're using a function
		// pointer hack, so use mutex
		enumWS.m.Lock()
		defer enumWS.m.Unlock()

		// Set global pointer to our local func
		enumWS.fp = enumFunc

		_, _, _ = user32.NewProc(proc).Call(
			uintptr(unsafe.Pointer(C.enumWindowStationsHook)),
			params,
		)

		// Set global pointer back to nil
		enumWS.fp = nil
	case uintptr:
		_, _, _ = user32.NewProc(proc).Call(enumFunc, params)
	default:
		return errors.Newf("invalid enumFunc type: %T", enumFunc)
	}

	return nil
}

// This is a hack to call a Go function even tho EnumWindowStationsW
// takes a C function pointer.
//
//export enumWindowStationsHook
func enumWindowStationsHook(name *uint16, params uintptr) bool {
	if enumWS.fp == nil {
		return false
	}

	return enumWS.fp(windows.UTF16PtrToString(name), params)
}

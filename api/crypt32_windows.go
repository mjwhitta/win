//go:build windows

package api

import (
	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var crypt32 *windows.LazyDLL = windows.NewLazySystemDLL("crypt32")

// CertEnumPhysicalStore is from wincrypt.h
func CertEnumPhysicalStore(
	store string,
	dwFlags uintptr,
	pvArg uintptr,
	pfnEnum uintptr,
) error {
	var e error
	var ok uintptr
	var proc string = "CertEnumPhysicalStore"

	ok, _, e = crypt32.NewProc(proc).Call(
		types.LpCwstr(store),
		dwFlags,
		pvArg,
		pfnEnum,
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// CertEnumSystemStore is from wincrypt.h
func CertEnumSystemStore(
	dwFlags uintptr,
	pvSystemStoreLocationPara uintptr,
	pvArg uintptr,
	pfnEnum uintptr,
) error {
	var e error
	var ok uintptr
	var proc string = "CertEnumSystemStore"

	ok, _, e = crypt32.NewProc(proc).Call(
		dwFlags,
		pvSystemStoreLocationPara,
		pvArg,
		pfnEnum,
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

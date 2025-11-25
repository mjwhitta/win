//go:build windows

package api

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/win/types"
)

var wininet *windows.LazyDLL = windows.NewLazySystemDLL("Wininet")

// HTTPAddRequestHeadersW from wininet.h
func HTTPAddRequestHeadersW(
	reqHndl uintptr,
	header string,
	addMethod uintptr,
) error {
	var e error
	var ok uintptr
	var proc string = "HttpAddRequestHeadersW"

	if header == "" {
		// Weird, just do nothing
		return nil
	}

	header = strings.TrimSpace(header) + "\r\n"

	ok, _, e = wininet.NewProc(proc).Call(
		reqHndl,
		types.LpCwstr(header),
		uintptr(len(header)),
		addMethod,
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// HTTPOpenRequestW from wininet.h
func HTTPOpenRequestW(
	connHndl uintptr,
	verb string,
	objectName string,
	version string,
	referrer string,
	acceptTypes []string,
	flags uintptr,
	context uintptr,
) (uintptr, error) {
	var e error
	var lplpcwstrAcceptTypes []*uint16
	var proc string = "HttpOpenRequestW"
	var reqHndl uintptr

	// Convert to Windows types
	for _, theType := range acceptTypes {
		if theType == "" {
			continue
		}

		lplpcwstrAcceptTypes = append(
			lplpcwstrAcceptTypes,
			types.Cwstr(theType),
		)
	}

	if len(lplpcwstrAcceptTypes) == 0 {
		lplpcwstrAcceptTypes = make([]*uint16, 1)
	}

	reqHndl, _, e = wininet.NewProc(proc).Call(
		connHndl,
		types.LpCwstr(verb),
		types.LpCwstr(objectName),
		types.LpCwstr(version),
		types.LpCwstr(referrer),
		uintptr(unsafe.Pointer(&lplpcwstrAcceptTypes[0])),
		flags,
		context,
	)
	if reqHndl == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return reqHndl, nil
}

// HTTPQueryInfoW from wininet.h
func HTTPQueryInfoW(
	reqHndl uintptr,
	info uintptr,
	buffer *[]byte,
	bufferLen *int,
	index *int,
) error {
	var b []uint16
	var e error
	var ok uintptr
	var proc string = "HttpQueryInfoW"
	var tmp string

	if *bufferLen > 0 {
		b = make([]uint16, *bufferLen)
	} else {
		b = make([]uint16, 1)
	}

	ok, _, e = wininet.NewProc(proc).Call(
		reqHndl,
		info,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(bufferLen)),
		uintptr(unsafe.Pointer(index)),
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	tmp = windows.UTF16ToString(b)
	*buffer = []byte(tmp)

	return nil
}

// HTTPSendRequestW from wininet.h
func HTTPSendRequestW(
	reqHndl uintptr,
	headers string,
	headersLen int,
	data []byte,
	dataLen int,
) error {
	var body uintptr
	var e error
	var ok uintptr
	var proc string = "HttpSendRequestW"

	// Pointer to data if provided
	if (data != nil) && (len(data) > 0) {
		body = uintptr(unsafe.Pointer(&data[0]))
	}

	ok, _, e = wininet.NewProc(proc).Call(
		reqHndl,
		types.LpCwstr(headers),
		uintptr(headersLen),
		body,
		uintptr(dataLen),
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// InternetCloseHandle from wininet.h
func InternetCloseHandle(reqHndl uintptr) error {
	var e error
	var ok uintptr
	var proc string = "InternetCloseHandle"

	if ok, _, e = wininet.NewProc(proc).Call(reqHndl); ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// InternetConnectW from wininet.h
func InternetConnectW(
	sessionHndl uintptr,
	serverName string,
	serverPort int,
	username string,
	password string,
	service uintptr,
	flags uintptr,
	context uintptr,
) (uintptr, error) {
	var connHndl uintptr
	var e error
	var proc string = "InternetConnectW"

	connHndl, _, e = wininet.NewProc(proc).Call(
		sessionHndl,
		types.LpCwstr(serverName),
		uintptr(serverPort),
		types.LpCwstr(username),
		types.LpCwstr(password),
		service,
		flags,
		context,
	)
	if connHndl == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return connHndl, nil
}

// InternetGetCookieW from wininet.h
func InternetGetCookieW(
	url string,
	buffer *[]byte,
	bufferLen *int,
) error {
	var b []uint16
	var e error
	var ok uintptr
	var proc string = "InternetGetCookieW"
	var tmp string

	if *bufferLen > 0 {
		b = make([]uint16, *bufferLen)
	} else {
		b = make([]uint16, 1)
	}

	ok, _, e = wininet.NewProc(proc).Call(
		types.LpCwstr(url),
		0,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(bufferLen)),
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	tmp = windows.UTF16ToString(b)
	*buffer = []byte(tmp)

	return nil
}

// InternetGetCookieExW from wininet.h
func InternetGetCookieExW(
	url string,
	name string,
	buffer *[]byte,
	bufferLen *int,
	flags uintptr,
) error {
	var b []uint16
	var e error
	var ok uintptr
	var proc string = "InternetGetCookieExW"
	var tmp string

	if *bufferLen > 0 {
		b = make([]uint16, *bufferLen)
	} else {
		b = make([]uint16, 1)
	}

	ok, _, e = wininet.NewProc(proc).Call(
		types.LpCwstr(url),
		types.LpCwstr(name),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(bufferLen)),
		flags,
		0,
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	tmp = windows.UTF16ToString(b)
	*buffer = []byte(tmp)

	return nil
}

// InternetOpenW from wininet.h
func InternetOpenW(
	userAgent string,
	accessType uintptr,
	proxy string,
	proxyBypass string,
	flags uintptr,
) (uintptr, error) {
	var e error
	var proc string = "InternetOpenW"
	var sessionHndl uintptr

	sessionHndl, _, e = wininet.NewProc(proc).Call(
		types.LpCwstr(userAgent),
		accessType,
		types.LpCwstr(proxy),
		types.LpCwstr(proxyBypass),
		flags,
	)
	if sessionHndl == 0 {
		return 0, errors.Newf("%s: %w", proc, e)
	}

	return sessionHndl, nil
}

// InternetQueryDataAvailable from wininet.h
func InternetQueryDataAvailable(
	reqHndl uintptr,
	bytesAvailable *int64,
) error {
	var e error
	var ok uintptr
	var proc string = "InternetQueryDataAvailable"

	ok, _, e = wininet.NewProc(proc).Call(
		reqHndl,
		uintptr(unsafe.Pointer(bytesAvailable)),
		0,
		0,
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

// InternetReadFile from wininet.h
func InternetReadFile(
	reqHndl uintptr,
	buffer *[]byte,
	bytesToRead int64,
	bytesRead *int64,
) error {
	var b []byte
	var e error
	var ok uintptr
	var proc string = "InternetReadFile"

	if bytesToRead > 0 {
		b = make([]byte, bytesToRead)
	} else {
		b = make([]byte, 1)
	}

	ok, _, e = wininet.NewProc(proc).Call(
		reqHndl,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(bytesToRead),
		uintptr(unsafe.Pointer(bytesRead)),
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	*buffer = b

	return nil
}

// InternetSetOptionW from wininet.h
func InternetSetOptionW(
	hndl uintptr,
	opt uintptr,
	val []byte,
	valLen int,
) error {
	var e error
	var ok uintptr
	var proc string = "InternetSetOptionW"

	// Pointer to data if provided
	if valLen == 0 {
		val = make([]byte, 1)
	}

	ok, _, e = wininet.NewProc(proc).Call(
		hndl,
		opt,
		uintptr(unsafe.Pointer(&val[0])),
		uintptr(valLen),
	)
	if ok == 0 {
		return errors.Newf("%s: %w", proc, e)
	}

	return nil
}

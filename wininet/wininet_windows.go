package wininet

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var wininet *windows.LazyDLL = windows.NewLazySystemDLL("Wininet")

// HTTPAddRequestHeadersW is from wininet.h
func HTTPAddRequestHeadersW(
	reqHndl uintptr,
	header string,
	addMethod uintptr,
) error {
	var e error
	var ok uintptr
	var pswzHeader uintptr
	var tmp *uint16

	if header == "" {
		// Weird, just do nothing
		return nil
	}

	header = strings.TrimSpace(header) + "\r\n"

	// Convert to Windows types
	if tmp, e = windows.UTF16PtrFromString(header); e != nil {
		return e
	}

	pswzHeader = uintptr(unsafe.Pointer(tmp))

	ok, _, e = wininet.NewProc("HttpAddRequestHeadersW").Call(
		reqHndl,
		pswzHeader,
		uintptr(len(header)),
		addMethod,
	)
	if ok == 0 {
		return fmt.Errorf("HttpAddRequestHeadersW: %s", e.Error())
	}

	return nil
}

// HTTPOpenRequestW is from wininet.h
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
	var lpcwstrObjectName uintptr
	var lpcwstrReferrer uintptr
	var lpcwstrVerb uintptr
	var lpcwstrVersion uintptr
	var lplpcwstrAcceptTypes []*uint16
	var reqHndl uintptr
	var tmp *uint16

	// Convert to Windows types
	lplpcwstrAcceptTypes = make([]*uint16, 1)
	for _, theType := range acceptTypes {
		if theType == "" {
			continue
		}

		tmp, e = windows.UTF16PtrFromString(theType)
		if e != nil {
			return 0, e
		}

		lplpcwstrAcceptTypes = append(lplpcwstrAcceptTypes, tmp)
	}

	if objectName != "" {
		if tmp, e = windows.UTF16PtrFromString(objectName); e != nil {
			return 0, e
		}

		lpcwstrObjectName = uintptr(unsafe.Pointer(tmp))
	}

	if referrer != "" {
		if tmp, e = windows.UTF16PtrFromString(referrer); e != nil {
			return 0, e
		}

		lpcwstrReferrer = uintptr(unsafe.Pointer(tmp))
	}

	if verb != "" {
		if tmp, e = windows.UTF16PtrFromString(verb); e != nil {
			return 0, e
		}

		lpcwstrVerb = uintptr(unsafe.Pointer(tmp))
	}

	if version != "" {
		if tmp, e = windows.UTF16PtrFromString(version); e != nil {
			return 0, e
		}

		lpcwstrVersion = uintptr(unsafe.Pointer(tmp))
	}

	reqHndl, _, e = wininet.NewProc("HttpOpenRequestW").Call(
		connHndl,
		lpcwstrVerb,
		lpcwstrObjectName,
		lpcwstrVersion,
		lpcwstrReferrer,
		uintptr(unsafe.Pointer(&lplpcwstrAcceptTypes[0])),
		flags,
		context,
	)
	if reqHndl == 0 {
		return 0, fmt.Errorf("HttpOpenRequest: %s", e.Error())
	}

	return reqHndl, nil
}

// HTTPQueryInfoW is from wininet.h
func HTTPQueryInfoW(
	reqHndl uintptr,
	info uintptr,
	buffer *[]byte,
	bufferLen *int,
	index *int,
) error {
	var b []uint16
	var e error
	var success uintptr
	var tmp string

	if *bufferLen > 0 {
		b = make([]uint16, *bufferLen)
	} else {
		b = make([]uint16, 1)
	}

	success, _, e = wininet.NewProc("HttpQueryInfoW").Call(
		reqHndl,
		info,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(bufferLen)),
		uintptr(unsafe.Pointer(index)),
	)
	if success == 0 {
		return fmt.Errorf("HttpQueryInfoW: %s", e.Error())
	}

	tmp = windows.UTF16ToString(b)
	*buffer = []byte(tmp)

	return nil
}

// HTTPSendRequestW is from wininet.h
func HTTPSendRequestW(
	reqHndl uintptr,
	headers string,
	headersLen int,
	data []byte,
	dataLen int,
) error {
	var body uintptr
	var e error
	var lpcwstrHeaders uintptr
	var success uintptr
	var tmp *uint16

	// Pointer to data if provided
	if (data != nil) && (len(data) > 0) {
		body = uintptr(unsafe.Pointer(&data[0]))
	}

	// Convert to Windows types
	if headersLen > 0 {
		if tmp, e = windows.UTF16PtrFromString(headers); e != nil {
			return e
		}

		lpcwstrHeaders = uintptr(unsafe.Pointer(tmp))
	}

	success, _, e = wininet.NewProc("HttpSendRequestW").Call(
		reqHndl,
		lpcwstrHeaders,
		uintptr(headersLen),
		body,
		uintptr(dataLen),
	)
	if success == 0 {
		return fmt.Errorf("HttpSendRequestW: %s", e.Error())
	}

	return nil
}

// InternetConnectW is from wininet.h
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
	var lpcwstrServerName uintptr
	var lpcwstrUserName uintptr
	var lpcwstrPassword uintptr
	var tmp *uint16

	// Convert to Windows types
	if password != "" {
		if tmp, e = windows.UTF16PtrFromString(password); e != nil {
			return 0, e
		}

		lpcwstrPassword = uintptr(unsafe.Pointer(tmp))
	}

	if serverName != "" {
		if tmp, e = windows.UTF16PtrFromString(serverName); e != nil {
			return 0, e
		}

		lpcwstrServerName = uintptr(unsafe.Pointer(tmp))
	}

	if username != "" {
		if tmp, e = windows.UTF16PtrFromString(username); e != nil {
			return 0, e
		}

		lpcwstrUserName = uintptr(unsafe.Pointer(tmp))
	}

	connHndl, _, e = wininet.NewProc("InternetConnectW").Call(
		sessionHndl,
		lpcwstrServerName,
		uintptr(serverPort),
		lpcwstrUserName,
		lpcwstrPassword,
		service,
		flags,
		context,
	)
	if connHndl == 0 {
		return 0, fmt.Errorf("InternetConnectW: %s", e.Error())
	}

	return connHndl, nil
}

// InternetOpenW is from wininet.h
func InternetOpenW(
	userAgent string,
	accessType uintptr,
	proxy string,
	proxyBypass string,
	flags uintptr,
) (uintptr, error) {
	var e error
	var lpszAgent uintptr
	var lpszProxy uintptr
	var lpszProxyBypass uintptr
	var sessionHndl uintptr
	var tmp *uint16

	// Convert to Windows types
	if userAgent != "" {
		if tmp, e = windows.UTF16PtrFromString(userAgent); e != nil {
			return 0, e
		}
		lpszAgent = uintptr(unsafe.Pointer(tmp))
	}

	if proxy != "" {
		if tmp, e = windows.UTF16PtrFromString(proxy); e != nil {
			return 0, e
		}
		lpszProxy = uintptr(unsafe.Pointer(tmp))
	}

	if proxyBypass != "" {
		tmp, e = windows.UTF16PtrFromString(proxyBypass)
		if e != nil {
			return 0, e
		}
		lpszProxyBypass = uintptr(unsafe.Pointer(tmp))
	}

	sessionHndl, _, e = wininet.NewProc("InternetOpenW").Call(
		lpszAgent,
		accessType,
		lpszProxy,
		lpszProxyBypass,
		flags,
	)
	if sessionHndl == 0 {
		return 0, fmt.Errorf("InternetOpenW: %s", e.Error())
	}

	return sessionHndl, nil
}

// InternetQueryDataAvailable is from wininet.h
func InternetQueryDataAvailable(
	reqHndl uintptr,
	bytesAvailable *int64,
) error {
	var e error
	var success uintptr

	success, _, e = wininet.NewProc(
		"InternetQueryDataAvailable",
	).Call(
		reqHndl,
		uintptr(unsafe.Pointer(bytesAvailable)),
		0,
		0,
	)
	if success == 0 {
		return fmt.Errorf("InternetQueryDataAvailable: %s", e.Error())
	}

	return nil
}

// InternetReadFile is from wininet.h
func InternetReadFile(
	reqHndl uintptr,
	buffer *[]byte,
	bytesToRead int64,
	bytesRead *int64,
) error {
	var b []byte
	var e error
	var success uintptr

	if bytesToRead > 0 {
		b = make([]byte, bytesToRead)
	} else {
		b = make([]byte, 1)
	}

	success, _, e = wininet.NewProc("InternetReadFile").Call(
		reqHndl,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(bytesToRead),
		uintptr(unsafe.Pointer(bytesRead)),
	)
	if success == 0 {
		return fmt.Errorf("InternetReadFile: %s", e.Error())
	}

	*buffer = b

	return nil
}

// InternetSetOptionW is from wininet.h
func InternetSetOptionW(
	hndl uintptr,
	opt uintptr,
	val []byte,
	valLen int,
) error {
	var e error
	var success uintptr

	// Pointer to data if provided
	if valLen == 0 {
		val = make([]byte, 1)
	}

	success, _, e = wininet.NewProc("InternetSetOptionW").Call(
		hndl,
		opt,
		uintptr(unsafe.Pointer(&val[0])),
		uintptr(valLen),
	)
	if success == 0 {
		return fmt.Errorf("InternetSetOptionW: %s", e.Error())
	}

	return nil
}

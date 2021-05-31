package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/mjwhitta/win/wininet"
)

func buildRequest(
	sessionHndl uintptr,
	method string,
	dst string,
) (uintptr, error) {
	var connHndl uintptr
	var e error
	var flags uintptr
	var passwd string
	var port int64
	var query string
	var reqHndl uintptr
	var uri *url.URL

	// Parse URL
	if uri, e = url.Parse(dst); e != nil {
		return 0, e
	}

	passwd, _ = uri.User.Password()

	if uri.Port() != "" {
		if port, e = strconv.ParseInt(uri.Port(), 10, 64); e != nil {
			return 0, e
		}
	}

	switch uri.Scheme {
	case "https":
		flags = wininet.InternetFlagSecure
	}

	// Create connection
	connHndl, e = wininet.InternetConnectW(
		sessionHndl,
		uri.Hostname(),
		int(port),
		uri.User.Username(),
		passwd,
		wininet.InternetServiceHTTP,
		flags,
		0,
	)
	if e != nil {
		return 0, e
	}

	// Send query string too
	if uri.RawQuery != "" {
		query = "?" + uri.RawQuery
	}

	// Create HTTP request
	reqHndl, e = wininet.HTTPOpenRequestW(
		connHndl,
		method,
		uri.Path+query,
		"",
		"",
		[]string{},
		flags,
		0,
	)
	if e != nil {
		return 0, e
	}

	return reqHndl, nil
}

func buildResponse(reqHndl uintptr) (*Response, error) {
	var b []byte
	var body io.ReadCloser
	var code int64
	var contentLen int64
	var e error
	var hdrs map[string][]string
	var major int
	var minor int
	var proto string
	var res *Response
	var status string

	// Get status code
	b, e = queryResponse(reqHndl, wininet.HTTPQueryStatusCode)
	if e != nil {
		return nil, e
	}

	status = string(b)
	if code, e = strconv.ParseInt(status, 10, 64); e != nil {
		return nil, e
	}

	// Get status text
	b, e = queryResponse(reqHndl, wininet.HTTPQueryStatusText)
	if e != nil {
		return nil, e
	} else if len(b) > 0 {
		status += " " + string(b)
	}

	// Parse headers and proto
	if proto, major, minor, hdrs, e = getHeaders(reqHndl); e != nil {
		return nil, e
	}

	// Get Content-Length
	b, e = queryResponse(reqHndl, wininet.HTTPQueryContentLength)
	if e == nil {
		contentLen, e = strconv.ParseInt(string(b), 10, 64)
		if e != nil {
			return nil, e
		}
	} else {
		contentLen = 0
	}

	// Read response body
	if body, e = readResponse(reqHndl, contentLen); e != nil {
		return nil, e
	}

	res = &Response{
		Body:          body,
		ContentLength: contentLen,
		Header:        hdrs,
		Proto:         proto,
		ProtoMajor:    major,
		ProtoMinor:    minor,
		Status:        status,
		StatusCode:    int(code),
	}

	return res, nil
}

func getHeaders(
	reqHndl uintptr,
) (string, int, int, map[string][]string, error) {
	var b []byte
	var e error
	var hdrs = map[string][]string{}
	var major int64
	var minor int64
	var proto string
	var tmp []string

	// Get headers
	b, e = queryResponse(reqHndl, wininet.HTTPQueryRawHeadersCRLF)
	if e != nil {
		return "", 0, 0, nil, e
	}

	for _, hdr := range strings.Split(string(b), "\r\n") {
		tmp = strings.SplitN(hdr, ": ", 2)

		if len(tmp) == 2 {
			if _, ok := hdrs[tmp[0]]; !ok {
				hdrs[tmp[0]] = []string{}
			}

			hdrs[tmp[0]] = append(hdrs[tmp[0]], tmp[1])
		} else if strings.HasPrefix(hdr, "HTTP") {
			proto = strings.Fields(hdr)[0]
			tmp = strings.Split(proto, ".")

			if len(tmp) >= 2 {
				tmp[0] = strings.Replace(tmp[0], "HTTP/", "", 1)

				major, e = strconv.ParseInt(tmp[0], 10, 64)
				if e != nil {
					return "", 0, 0, nil, e
				}

				minor, e = strconv.ParseInt(tmp[1], 10, 64)
				if e != nil {
					return "", 0, 0, nil, e
				}
			}
		}
	}

	return proto, int(major), int(minor), hdrs, nil
}

func queryResponse(reqHndl uintptr, info uintptr) ([]byte, error) {
	var buffer []byte
	var e error
	var size int

	e = wininet.HTTPQueryInfoW(reqHndl, info, &buffer, &size, 0)
	if e != nil {
		buffer = make([]byte, size)

		e = wininet.HTTPQueryInfoW(reqHndl, info, &buffer, &size, 0)
		if e != nil {
			return []byte{}, e
		}
	}

	return buffer, nil
}

func readResponse(
	reqHndl uintptr,
	contentLen int64,
) (io.ReadCloser, error) {
	var b []byte
	var n int64
	var e error

	e = wininet.InternetReadFile(reqHndl, &b, contentLen, &n)
	if e != nil {
		return nil, e
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func sendRequest(
	reqHndl uintptr,
	headers map[string]string,
	data []byte,
) error {
	var combinedHdrs string

	// Combine headers
	if headers != nil {
		for k, v := range headers {
			combinedHdrs += "\n\r" + k + ": " + v
		}
		combinedHdrs = strings.TrimSpace(combinedHdrs)
	}

	// Send HTTP request
	return wininet.HTTPSendRequestW(
		reqHndl,
		combinedHdrs,
		len([]byte(combinedHdrs)),
		data,
		len(data),
	)
}

package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/mjwhitta/win/winhttp"
)

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

	// Get response
	if e = winhttp.ReceiveResponse(reqHndl); e != nil {
		return nil, e
	}

	// Get status code
	b, e = queryResponse(reqHndl, winhttp.WinhttpQueryStatusCode)
	if e != nil {
		return nil, e
	}

	status = string(b)
	if code, e = strconv.ParseInt(status, 10, 64); e != nil {
		return nil, e
	}

	// Get status text
	b, e = queryResponse(reqHndl, winhttp.WinhttpQueryStatusText)
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
	b, e = queryResponse(reqHndl, winhttp.WinhttpQueryContentLength)
	if e != nil {
		return nil, e
	}

	if contentLen, e = strconv.ParseInt(string(b), 10, 64); e != nil {
		return nil, e
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
	b, e = queryResponse(reqHndl, winhttp.WinhttpQueryRawHeadersCRLF)
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

	e = winhttp.QueryHeaders(reqHndl, info, "", &buffer, &size, 0)
	if e != nil {
		buffer = make([]byte, size)

		e = winhttp.QueryHeaders(reqHndl, info, "", &buffer, &size, 0)
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

	if e = winhttp.ReadData(reqHndl, &b, contentLen, &n); e != nil {
		return nil, e
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}

func sendRequest(
	sessionHndl uintptr,
	method string,
	dst string,
	headers map[string]string,
	data []byte,
) (uintptr, error) {
	var combinedHdrs string
	var connHndl uintptr
	var e error
	var flags uintptr
	var port int64
	var reqHndl uintptr
	var uri *url.URL

	// Parse URL
	if uri, e = url.Parse(dst); e != nil {
		return 0, e
	}

	if uri.Port() != "" {
		if port, e = strconv.ParseInt(uri.Port(), 10, 64); e != nil {
			return 0, e
		}
	}

	switch uri.Scheme {
	case "https":
		flags = winhttp.WinhttpFlagSecure
	}

	// Combine headers
	if headers != nil {
		for k, v := range headers {
			combinedHdrs += "\n\r" + k + ": " + v
		}
		combinedHdrs = strings.TrimSpace(combinedHdrs)
	}

	// Create connection
	connHndl, e = winhttp.Connect(
		sessionHndl,
		uri.Hostname(),
		int(port),
	)
	if e != nil {
		return 0, e
	}

	// Create HTTP request
	reqHndl, e = winhttp.OpenRequest(
		connHndl,
		method,
		uri.Path,
		"",
		"",
		"",
		flags,
	)
	if e != nil {
		return 0, e
	}

	// Send HTTP request
	e = winhttp.SendRequest(
		reqHndl,
		combinedHdrs,
		len([]byte(combinedHdrs)),
		data,
		len(data),
	)
	if e != nil {
		return 0, e
	}

	return reqHndl, nil
}

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

func buildRequest(sessionHndl uintptr, r *Request) (uintptr, error) {
	var connHndl uintptr
	var e error
	var flags uintptr
	var passwd string
	var port int64
	var query string
	var reqHndl uintptr
	var uri *url.URL

	// Parse URL
	if uri, e = url.Parse(r.URL); e != nil {
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

	// Allow NTLM auth
	flags |= wininet.InternetFlagKeepConnection

	// Create HTTP request
	reqHndl, e = wininet.HTTPOpenRequestW(
		connHndl,
		r.Method,
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

var cookies []*Cookie

func buildResponse(reqHndl uintptr, req *Request) (*Response, error) {
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
	b, e = queryResponse(reqHndl, wininet.HTTPQueryStatusCode, 0)
	if e != nil {
		return nil, e
	}

	status = string(b)
	if code, e = strconv.ParseInt(status, 10, 64); e != nil {
		return nil, e
	}

	// Get status text
	b, e = queryResponse(reqHndl, wininet.HTTPQueryStatusText, 0)
	if e != nil {
		return nil, e
	} else if len(b) > 0 {
		status += " " + string(b)
	}

	// Parse cookies
	if cookies, e = getCookies(reqHndl); e != nil {
		return nil, e
	}

	// Parse headers and proto
	if proto, major, minor, hdrs, e = getHeaders(reqHndl); e != nil {
		return nil, e
	}

	// Read response body
	if body, contentLen, e = readResponse(reqHndl); e != nil {
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

	// Concat all cookies
	for _, c := range req.Cookies() {
		res.AddCookie(c)
	}

	for _, c := range cookies {
		res.AddCookie(c)
	}

	return res, nil
}

func getCookies(reqHndl uintptr) ([]*Cookie, error) {
	var b []byte
	var cookies []*Cookie
	var e error
	var tmp []string

	// Get cookies
	for i := 0; e == nil; i++ {
		b, e = queryResponse(
			reqHndl,
			wininet.HTTPQuerySetCookie,
			i,
		)
		if e == nil {
			tmp = strings.SplitN(string(b), "=", 2)
			cookies = append(
				cookies,
				&Cookie{Name: tmp[0], Value: tmp[1]},
			)
		}
	}

	return cookies, nil
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
	b, e = queryResponse(reqHndl, wininet.HTTPQueryRawHeadersCRLF, 0)
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

func queryResponse(reqHndl, info uintptr, idx int) ([]byte, error) {
	var buffer []byte
	var e error
	var size int

	if idx < 0 {
		idx = 0
	}

	e = wininet.HTTPQueryInfoW(reqHndl, info, &buffer, &size, &idx)
	if e != nil {
		buffer = make([]byte, size)

		e = wininet.HTTPQueryInfoW(
			reqHndl,
			info,
			&buffer,
			&size,
			&idx,
		)
		if e != nil {
			return []byte{}, e
		}
	}

	return buffer, nil
}

func readResponse(reqHndl uintptr) (io.ReadCloser, int64, error) {
	var b []byte
	var chunk []byte
	var chunkLen int64
	var contentLen int64
	var e error
	var n int64

	// Get Content-Length and body of response
	for {
		// Get next chunk size
		e = wininet.InternetQueryDataAvailable(reqHndl, &chunkLen)
		if e != nil {
			break
		}

		// Stop, if finished
		if chunkLen == 0 {
			break
		}

		// Read next chunk
		e = wininet.InternetReadFile(reqHndl, &chunk, chunkLen, &n)
		if e != nil {
			break
		}

		// Update fields
		contentLen += chunkLen
		b = append(b, chunk...)
	}

	if e != nil {
		return nil, 0, e
	}

	return ioutil.NopCloser(bytes.NewReader(b)), contentLen, nil
}

func sendRequest(reqHndl uintptr, r *Request) error {
	var e error
	var method uintptr

	// Process cookies
	method = wininet.HTTPAddreqFlagAdd
	// FIXME why doesn't this work here?!
	// method |= wininet.HTTPAddreqFlagCoalesceWithSemicolon

	// FIXME This is a dumb hack
	wininet.HTTPAddRequestHeadersW(
		reqHndl,
		"Cookie: ignore=ignore",
		wininet.HTTPAddreqFlagAddIfNew,
	)
	// End dumb hack

	for _, c := range r.Cookies() {
		e = wininet.HTTPAddRequestHeadersW(
			reqHndl,
			"Cookie: "+c.Name+"="+c.Value,
			method,
		)
		if e != nil {
			return e
		}
	}

	// Process headers
	method = wininet.HTTPAddreqFlagAdd
	method |= wininet.HTTPAddreqFlagReplace

	for k, v := range r.Headers {
		e = wininet.HTTPAddRequestHeadersW(
			reqHndl,
			k+": "+v,
			method,
		)
		if e != nil {
			return e
		}
	}

	// Send HTTP request
	return wininet.HTTPSendRequestW(
		reqHndl,
		"",
		0,
		r.Body,
		len(r.Body),
	)
}

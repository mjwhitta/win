//go:build windows

package winhttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

func buildRequest(
	sessionHndl uintptr, req *Request,
) (uintptr, error) {
	var connHndl uintptr
	var e error
	var flags uintptr
	var port int64
	var query string
	var reqHndl uintptr
	var uri *url.URL

	// Parse URL
	if uri, e = url.Parse(req.URL); e != nil {
		e = errors.Newf("failed to parse url %s: %w", req.URL, e)
		return 0, e
	}

	if uri.Port() != "" {
		if port, e = strconv.ParseInt(uri.Port(), 10, 64); e != nil {
			e = errors.Newf("port %s invalid: %w", uri.Port(), e)
			return 0, e
		}
	}

	switch uri.Scheme {
	case "https":
		flags = w32.Winhttp.WinhttpFlagSecure
	}

	// Create connection
	connHndl, e = w32.WinHTTPConnect(
		sessionHndl,
		uri.Hostname(),
		int(port),
	)
	if e != nil {
		return 0, errors.Newf("failed to create connection: %w", e)
	}

	// Send query string too
	if uri.RawQuery != "" {
		query = "?" + uri.RawQuery
	}

	// Create HTTP request
	reqHndl, e = w32.WinHTTPOpenRequest(
		connHndl,
		req.Method,
		uri.Path+query,
		"",
		"",
		[]string{},
		flags,
	)
	if e != nil {
		return 0, errors.Newf("failed to open request: %w", e)
	}

	return reqHndl, nil
}

func buildResponse(reqHndl uintptr, req *Request) (*Response, error) {
	var b []byte
	var body io.ReadCloser
	var code int64
	var contentLen int64
	var cookies []*Cookie
	var e error
	var hdrs map[string][]string
	var major int
	var minor int
	var proto string
	var res *Response
	var status string

	// Get response
	if e = w32.WinHTTPReceiveResponse(reqHndl); e != nil {
		return nil, errors.Newf("failed to get response: %w", e)
	}

	// Get status code
	b, e = queryResponse(
		reqHndl,
		w32.Winhttp.WinhttpQueryStatusCode,
		0,
	)
	if e != nil {
		return nil, e
	}

	status = string(b)
	if code, e = strconv.ParseInt(status, 10, 64); e != nil {
		return nil, errors.Newf("status %s invalid: %w", status, e)
	}

	// Get status text
	b, e = queryResponse(
		reqHndl,
		w32.Winhttp.WinhttpQueryStatusText,
		0,
	)
	if e != nil {
		return nil, e
	} else if len(b) > 0 {
		status += " " + string(b)
	}

	// Parse cookies
	cookies = getCookies(reqHndl)

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

func getCookies(reqHndl uintptr) []*Cookie {
	var b []byte
	var cookies []*Cookie
	var e error

	// Get cookies
	for i := 0; ; i++ {
		b, e = queryResponse(
			reqHndl,
			w32.Winhttp.WinhttpQuerySetCookie,
			i,
		)
		if e != nil {
			break
		}

		cookies = append(cookies, &Cookie{parseCookie(string(b))})
	}

	return cookies
}

func getHeaders(
	reqHndl uintptr,
) (string, int, int, map[string][]string, error) {
	var b []byte
	var e error
	var hdrs map[string][]string = map[string][]string{}
	var major int64
	var minor int64
	var proto string
	var tmp []string

	// Get headers
	b, e = queryResponse(
		reqHndl,
		w32.Winhttp.WinhttpQueryRawHeadersCRLF,
		0,
	)
	if e != nil {
		return "", 0, 0, nil, e
	}

	for _, hdr := range strings.Split(string(b), "\req\n") {
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
					e = errors.Newf("invalid HTTP version: %w", e)
					return "", 0, 0, nil, e
				}

				minor, e = strconv.ParseInt(tmp[1], 10, 64)
				if e != nil {
					e = errors.Newf("invalid HTTP version: %w", e)
					return "", 0, 0, nil, e
				}
			}
		}
	}

	return proto, int(major), int(minor), hdrs, nil
}

func parseCookie(raw string) http.Cookie {
	var hdr http.Header = http.Header{}
	var req *http.Request

	hdr.Add("Cookie", raw)
	req = &http.Request{Header: hdr}

	return *req.Cookies()[0]
}

func queryResponse(reqHndl, info uintptr, idx int) ([]byte, error) {
	var buffer []byte
	var e error
	var size int

	if idx < 0 {
		idx = 0
	}

	e = w32.WinHTTPQueryHeaders(
		reqHndl,
		info,
		"",
		&buffer,
		&size,
		&idx,
	)
	if e != nil {
		buffer = make([]byte, size)

		e = w32.WinHTTPQueryHeaders(
			reqHndl,
			info,
			"",
			&buffer,
			&size,
			&idx,
		)
		if e != nil {
			e = errors.Newf("failed to query info: %w", e)
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
		e = w32.WinHTTPQueryDataAvailable(reqHndl, &chunkLen)
		if e != nil {
			e = errors.Newf("failed to query data available: %w", e)
			break
		}

		// Stop, if finished
		if chunkLen == 0 {
			break
		}

		// Read next chunk
		e = w32.WinHTTPReadData(reqHndl, &chunk, chunkLen, &n)
		if e != nil {
			e = errors.Newf("failed to read data: %w", e)
			break
		}

		// Update fields
		contentLen += chunkLen
		b = append(b, chunk...)
	}

	if e != nil {
		return nil, 0, e
	}

	return io.NopCloser(bytes.NewReader(b)), contentLen, nil
}

func retrieveCookies(jar http.CookieJar, uri *url.URL) []*Cookie {
	var tmp []*Cookie

	if jar == nil {
		return nil
	}

	for _, cookie := range jar.Cookies(uri) {
		tmp = append(tmp, &Cookie{*cookie})
	}

	return tmp
}

func sendRequest(reqHndl uintptr, req *Request) error {
	var e error
	var method uintptr

	// Process cookies
	method = w32.Winhttp.WinhttpAddreqFlagAdd
	method |= w32.Winhttp.WinhttpAddreqFlagCoalesceWithSemicolon

	for _, c := range req.Cookies() {
		e = w32.WinHTTPAddRequestHeaders(
			reqHndl,
			"Cookie: "+c.Name+"="+c.Value,
			method,
		)
		if e != nil {
			return errors.Newf("failed to add cookies: %w", e)
		}
	}

	// Process headers
	method = w32.Winhttp.WinhttpAddreqFlagAdd
	method |= w32.Winhttp.WinhttpAddreqFlagReplace

	for k, v := range req.Headers {
		e = w32.WinHTTPAddRequestHeaders(
			reqHndl,
			k+": "+v,
			method,
		)
		if e != nil {
			return errors.Newf("failed to add request headers: %w", e)
		}
	}

	// Send HTTP request
	e = w32.WinHTTPSendRequest(
		reqHndl,
		"",
		0,
		req.Body,
		len(req.Body),
	)
	if e != nil {
		return errors.Newf("failed to send request: %w", e)
	}

	return nil
}

func storeCookies(
	jar http.CookieJar, uri *url.URL, cookies []*Cookie,
) {
	var tmp []*http.Cookie

	if jar == nil {
		return
	}

	for _, cookie := range cookies {
		tmp = append(
			tmp,
			&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			},
		)
	}

	jar.SetCookies(uri, tmp)
}

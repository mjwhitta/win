//go:build windows

package winhttp

import (
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

func buildRequest(
	sessionHndl uintptr, req *http.Request, timeout time.Duration,
) (uintptr, error) {
	var b []byte = make([]byte, 4)
	var connHndl uintptr
	var e error
	var flags uintptr
	var port int64
	var query string
	var reqHndl uintptr

	if req.URL.Port() != "" {
		// If invalid port, Port() returns empty string, so no errors
		port, _ = strconv.ParseInt(req.URL.Port(), 10, 64)
	}

	switch req.URL.Scheme {
	case "https":
		flags = w32.Winhttp.WinhttpFlagSecure
	}

	// Create connection
	connHndl, e = w32.WinHTTPConnect(
		sessionHndl,
		req.URL.Hostname(),
		int(port),
	)
	if e != nil {
		return 0, errors.Newf("failed to create connection: %w", e)
	}

	// Send query string too
	if req.URL.RawQuery != "" {
		query = "?" + req.URL.RawQuery
	}

	// Send fragment too
	if req.URL.RawFragment != "" {
		query = "#" + req.URL.RawFragment
	}

	// Create HTTP request
	reqHndl, e = w32.WinHTTPOpenRequest(
		connHndl,
		req.Method,
		req.URL.Path+query,
		"",
		"",
		[]string{},
		flags,
	)
	if e != nil {
		return 0, errors.Newf("failed to open request: %w", e)
	}

	// Don't redirect
	flags = w32.Winhttp.WinhttpDisableRedirects

	// Don't let Windows handle cookies
	flags |= w32.Winhttp.WinhttpDisableCookies

	binary.LittleEndian.PutUint32(b, uint32(flags))
	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionDisableFeature,
		b,
		len(b),
	)
	if e != nil {
		e = errors.Newf("failed to set options: %w", e)
		return 0, e
	}

	if e = setTimeouts(reqHndl, timeout); e != nil {
		return 0, e
	}

	return reqHndl, nil
}

func buildResponse(
	reqHndl uintptr, req *http.Request,
) (*http.Response, error) {
	var b []byte
	var body io.ReadCloser
	var code int64
	var contentLen int64
	var e error
	var hdrs http.Header
	var major int
	var minor int
	var proto string
	var res *http.Response
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

	// Parse headers and proto
	if proto, major, minor, hdrs, e = getHeaders(reqHndl); e != nil {
		return nil, e
	}

	// Read response body
	if body, contentLen, e = readResponse(reqHndl); e != nil {
		return nil, e
	}

	res = &http.Response{
		Body:          body,
		ContentLength: contentLen,
		Header:        hdrs,
		Proto:         proto,
		ProtoMajor:    major,
		ProtoMinor:    minor,
		Request:       req,
		Status:        status,
		StatusCode:    int(code),
	}

	return res, nil
}

func dbgLog(debug bool, thing any) {
	var b []byte
	var e error

	if !debug {
		return
	}

	switch thing := thing.(type) {
	case *http.Request:
		if b, e = httputil.DumpRequestOut(thing, true); e == nil {
			println(string(b))
		}
	case *http.Response:
		if b, e = httputil.DumpResponse(thing, true); e == nil {
			println(string(b))
		}
	default:
		println(thing)
	}
}

func disableTLS(reqHndl uintptr) error {
	var b []byte = make([]byte, 4)
	var e error
	var tmp uintptr

	tmp |= w32.Winhttp.SecurityFlagIgnoreUnknownCa
	tmp |= w32.Winhttp.SecurityFlagIgnoreCertDateInvalid
	tmp |= w32.Winhttp.SecurityFlagIgnoreCertCnInvalid
	tmp |= w32.Winhttp.SecurityFlagIgnoreCertWrongUsage

	binary.LittleEndian.PutUint32(b, uint32(tmp))

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionSecurityFlags,
		b,
		len(b),
	)
	if e != nil {
		e = errors.Newf("failed to disable TLS verification: %w", e)
		return e
	}

	return nil
}

func getHeaders(
	reqHndl uintptr,
) (string, int, int, http.Header, error) {
	var b []byte
	var e error
	var hdrs http.Header = http.Header{}
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

func loadCookies(jar http.CookieJar, req *http.Request) {
	if jar == nil {
		return
	}

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}
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
		if size == 0 {
			return nil, errors.Newf("failed to query info: %w", e)
		}

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
			return nil, errors.Newf("failed to query info: %w", e)
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

func sendRequest(
	reqHndl uintptr, req *http.Request,
) (*http.Response, error) {
	var b []byte
	var e error
	var method uintptr
	var res *http.Response

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
			return nil, errors.Newf("failed to add cookies: %w", e)
		}
	}

	// Process headers
	method = w32.Winhttp.WinhttpAddreqFlagAdd
	method |= w32.Winhttp.WinhttpAddreqFlagReplace

	for k := range req.Header {
		e = w32.WinHTTPAddRequestHeaders(
			reqHndl,
			k+": "+req.Header.Get(k),
			method,
		)
		if e != nil {
			e = errors.Newf("failed to add request headers: %w", e)
			return nil, e
		}
	}

	if req.Body != nil {
		if b, e = io.ReadAll(req.Body); e != nil {
			e = errors.Newf("failed to read request body: %w", e)
			return nil, e
		}

		req.Body.Close()
	}

	// Send HTTP request
	e = w32.WinHTTPSendRequest(
		reqHndl,
		"",
		0,
		b,
		len(b),
	)
	if e != nil {
		e = errors.Newf("%s \"%s\": %w", req.Method, req.URL, e)
		return nil, e
	}

	if res, e = buildResponse(reqHndl, req); e != nil {
		return nil, e
	}

	return res, nil
}

func setTimeouts(reqHndl uintptr, timeout time.Duration) error {
	var b []byte
	var e error

	if timeout <= 0 {
		return nil
	}

	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(timeout.Milliseconds()))

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionConnectTimeout,
		b,
		len(b),
	)
	if e != nil {
		return errors.Newf("failed to set connect timeout: %w", e)
	}

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionReceiveResponseTimeout,
		b,
		len(b),
	)
	if e != nil {
		return errors.Newf("failed to set response timeout: %w", e)
	}

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionReceiveTimeout,
		b,
		len(b),
	)
	if e != nil {
		return errors.Newf("failed to set receive timeout: %w", e)
	}

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionResolveTimeout,
		b,
		len(b),
	)
	if e != nil {
		return errors.Newf("failed to set resolve timeout: %w", e)
	}

	e = w32.WinHTTPSetOption(
		reqHndl,
		w32.Winhttp.WinhttpOptionSendTimeout,
		b,
		len(b),
	)
	if e != nil {
		return errors.Newf("failed to set send timeout: %w", e)
	}

	return nil
}

func storeCookies(
	jar http.CookieJar, uri *url.URL, cookies []*http.Cookie,
) error {
	var e error
	var path *url.URL

	if jar == nil {
		return nil
	}

	// Store cookies per path
	for _, cookie := range cookies {
		if path, e = uri.Parse(cookie.Path); e != nil {
			return errors.Newf("invalid cookie path: %w", e)
		}

		jar.SetCookies(path, []*http.Cookie{cookie})
	}

	return nil
}

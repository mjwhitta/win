package winhttp

import (
	"encoding/binary"
	"time"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

// Client is a struct containing relevant metadata to make HTTP
// requests.
type Client struct {
	hndl            uintptr
	Timeout         time.Duration
	TLSClientConfig struct {
		InsecureSkipVerify bool
	}
}

// NewClient will return a pointer to a new Client instance that
// simply wraps the net/http.Client type.
func NewClient() (*Client, error) {
	var c = &Client{}
	var e error

	// Create session
	c.hndl, e = w32.WinHTTPOpen(
		"Go-http-client/1.1",
		w32.WinhttpAccessTypeAutomaticProxy,
		"",
		"",
		0,
	)
	if e != nil {
		return nil, errors.Newf("failed to create session: %w", e)
	}

	return c, nil
}

// Do will send the HTTP request and return an HTTP response.
func (c *Client) Do(r *Request) (*Response, error) {
	var b []byte
	var e error
	var reqHndl uintptr
	var res *Response
	var tlsIgnore uintptr

	if reqHndl, e = buildRequest(c.hndl, r); e != nil {
		return nil, e
	}

	if c.Timeout > 0 {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			b,
			uint32(c.Timeout.Milliseconds()),
		)

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionConnectTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set connect timeout: %w", e)
			return nil, e
		}

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionReceiveResponseTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set response timeout: %w", e)
			return nil, e
		}

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionReceiveTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set receive timeout: %w", e)
			return nil, e
		}

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionResolveTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set resolve timeout: %w", e)
			return nil, e
		}

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionSendTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set send timeout: %w", e)
			return nil, e
		}
	}

	if c.TLSClientConfig.InsecureSkipVerify {
		tlsIgnore |= w32.SecurityFlagIgnoreUnknownCa
		tlsIgnore |= w32.SecurityFlagIgnoreCertDateInvalid
		tlsIgnore |= w32.SecurityFlagIgnoreCertCnInvalid
		tlsIgnore |= w32.SecurityFlagIgnoreCertWrongUsage

		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(tlsIgnore))

		e = w32.WinHTTPSetOption(
			reqHndl,
			w32.WinhttpOptionSecurityFlags,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set security flags: %w", e)
			return nil, e
		}
	}

	if e = sendRequest(reqHndl, r); e != nil {
		return nil, e
	}

	if res, e = buildResponse(reqHndl, r); e != nil {
		return nil, e
	}

	return res, nil
}

// Get will make a GET request using w32.WinHTTPdll.
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(NewRequest(MethodGet, url))
}

// Head will make a HEAD request using w32.WinHTTPdll.
func (c *Client) Head(url string) (*Response, error) {
	return c.Do(NewRequest(MethodHead, url))
}

// Post will make a POST request using w32.WinHTTPdll.
func (c *Client) Post(
	url string,
	contentType string,
	body []byte,
) (*Response, error) {
	var r *Request = NewRequest(MethodPost, url, body)

	if contentType != "" {
		r.Headers["Content-Type"] = contentType
	}

	return c.Do(r)
}

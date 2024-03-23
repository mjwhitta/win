//go:build windows

package wininet

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
	var c *Client = &Client{}
	var e error

	// Create session
	c.hndl, e = w32.InternetOpenW(
		"Go-http-client/1.1", // TODO make this configurable
		w32.Wininet.InternetOpenTypePreconfig,
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

	if reqHndl, e = buildRequest(c.hndl, r); e != nil {
		return nil, e
	}

	if c.Timeout > 0 {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			b,
			uint32(c.Timeout.Milliseconds()),
		)

		e = w32.InternetSetOptionW(
			reqHndl,
			w32.Wininet.InternetOptionConnectTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set connect timeout: %w", e)
			return nil, e
		}

		e = w32.InternetSetOptionW(
			reqHndl,
			w32.Wininet.InternetOptionReceiveTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set receive timeout: %w", e)
			return nil, e
		}

		e = w32.InternetSetOptionW(
			reqHndl,
			w32.Wininet.InternetOptionSendTimeout,
			b,
			len(b),
		)
		if e != nil {
			e = errors.Newf("failed to set send timeout: %w", e)
			return nil, e
		}
	}

	if c.TLSClientConfig.InsecureSkipVerify {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			b,
			uint32(w32.Wininet.SecuritySetMask),
		)

		e = w32.InternetSetOptionW(
			reqHndl,
			w32.Wininet.InternetOptionSecurityFlags,
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

// Get will make a GET request using WinINet.dll.
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(NewRequest(MethodGet, url))
}

// Head will make a HEAD request using WinINet.dll.
func (c *Client) Head(url string) (*Response, error) {
	return c.Do(NewRequest(MethodHead, url))
}

// Post will make a POST request using WinINet.dll.
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

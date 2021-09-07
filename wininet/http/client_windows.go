package http

import (
	"encoding/binary"

	"gitlab.com/mjwhitta/win/wininet"
)

// NewClient will return a pointer to a new Client instnace that
// simply wraps the net/http.Client type.
func NewClient() (*Client, error) {
	var c = &Client{}
	var e error

	// Create session
	c.hndl, e = wininet.InternetOpenW(
		"Go-http-client/1.1",
		wininet.InternetOpenTypePreconfig,
		"",
		"",
		0,
	)
	if e != nil {
		return nil, e
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

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionConnectTimeout,
			b,
			len(b),
		)
		if e != nil {
			return nil, e
		}

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionReceiveTimeout,
			b,
			len(b),
		)
		if e != nil {
			return nil, e
		}

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionSendTimeout,
			b,
			len(b),
		)
		if e != nil {
			return nil, e
		}
	}

	if c.TLSClientConfig.InsecureSkipVerify {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			b,
			uint32(wininet.SecuritySetMask),
		)

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionSecurityFlags,
			b,
			len(b),
		)
		if e != nil {
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

// Get will make a GET request using Wininet.dll.
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(NewRequest(MethodGet, url))
}

// Head will make a HEAD request using Wininet.dll.
func (c *Client) Head(url string) (*Response, error) {
	return c.Do(NewRequest(MethodHead, url))
}

// Post will make a POST request using Wininet.dll.
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

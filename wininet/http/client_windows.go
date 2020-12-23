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

// Get will make a GET request using Wininet.dll.
func (c *Client) Get(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodGet, dst, headers, nil)
}

// Head will make a HEAD request using Wininet.dll.
func (c *Client) Head(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodHead, dst, headers, nil)
}

// Post will make a POST request using Wininet.dll.
func (c *Client) Post(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return c.request(MethodPost, dst, headers, data)
}

// Put will make a PUT request using Wininet.dll.
func (c *Client) Put(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return c.request(MethodPut, dst, headers, data)
}

func (c *Client) request(
	method string,
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	var e error
	var reqHndl uintptr
	var res *Response
	var tmp []byte

	if reqHndl, e = buildRequest(c.hndl, method, dst); e != nil {
		return nil, e
	}

	if c.Timeout > 0 {
		tmp = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			tmp,
			uint32(c.Timeout.Milliseconds()),
		)

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionConnectTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionReceiveTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionSendTimeout,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}
	}

	if c.TLSClientConfig.InsecureSkipVerify {
		tmp = make([]byte, 4)
		binary.LittleEndian.PutUint32(
			tmp,
			uint32(wininet.SecuritySetMask),
		)

		e = wininet.InternetSetOptionW(
			reqHndl,
			wininet.InternetOptionSecurityFlags,
			tmp,
			len(tmp),
		)
		if e != nil {
			return nil, e
		}
	}

	if e = sendRequest(reqHndl, headers, data); e != nil {
		return nil, e
	}

	if res, e = buildResponse(reqHndl); e != nil {
		return nil, e
	}

	return res, nil
}

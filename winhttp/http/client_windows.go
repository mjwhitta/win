package http

import "gitlab.com/mjwhitta/win/winhttp"

// NewClient will return a pointer to a new Client instnace that
// simply wraps the net/http.Client type.
func NewClient() (*Client, error) {
	var c = &Client{}
	var e error

	// Create session
	c.hndl, e = winhttp.Open(
		"Go-http-client/1.1",
		winhttp.WinhttpAccessTypeAutomaticProxy,
		"",
		"",
		0,
	)
	if e != nil {
		return nil, e
	}

	return c, nil
}

// Get will make a GET request using Winhttp.dll.
func (c *Client) Get(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodGet, dst, headers, nil)
}

// Head will make a HEAD request using Winhttp.dll.
func (c *Client) Head(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return c.request(MethodHead, dst, headers, nil)
}

// Post will make a POST request using Winhttp.dll.
func (c *Client) Post(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return c.request(MethodPost, dst, headers, data)
}

// Put will make a PUT request using Winhttp.dll.
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

	reqHndl, e = sendRequest(c.hndl, method, dst, headers, data)
	if e != nil {
		return nil, e
	}

	if res, e = buildResponse(reqHndl); e != nil {
		return nil, e
	}

	return res, nil
}

//go:build windows

package wininet

import (
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

// Client is a struct containing relevant metadata to make HTTP
// requests.
type Client struct {
	Jar     http.CookieJar
	Timeout time.Duration

	hndl uintptr
}

// NewClient will return a pointer to a new Client instance that
// simply wraps the net/http.Client type.
func NewClient(ua ...string) (*Client, error) {
	var c *Client = &Client{}
	var e error

	if len(ua) == 0 {
		ua = []string{"Go-http-client/1.1"}
	}

	// Create session
	c.hndl, e = w32.InternetOpenW(
		strings.Join(ua, " "),
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
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var b []byte
	var e error
	var reqHndl uintptr
	var res *http.Response

	for _, cookie := range loadCookies(c.Jar, req.URL) {
		req.AddCookie(cookie)
	}

	if reqHndl, e = buildRequest(c.hndl, req); e != nil {
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

	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		if t.TLSClientConfig != nil {
			if t.TLSClientConfig.InsecureSkipVerify {
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
					return nil, errors.Newf(
						"failed to disable TLS verification: %w",
						e,
					)
				}
			}
		}
	}

	if e = sendRequest(reqHndl, req); e != nil {
		return nil, e
	}

	if res, e = buildResponse(reqHndl, req); e != nil {
		return nil, e
	}

	storeCookies(c.Jar, req.URL, res.Cookies())

	return res, nil
}

// Get will make a GET request using WinINet.dll.
func (c *Client) Get(url string) (*http.Response, error) {
	var e error
	var req *http.Request

	if req, e = http.NewRequest(http.MethodGet, url, nil); e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	return c.Do(req)
}

// Head will make a HEAD request using WinINet.dll.
func (c *Client) Head(url string) (*http.Response, error) {
	var e error
	var req *http.Request

	if req, e = http.NewRequest(http.MethodHead, url, nil); e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	return c.Do(req)
}

// Post will make a POST request using WinINet.dll.
func (c *Client) Post(
	url string, contentType string, body io.Reader,
) (*http.Response, error) {
	var e error
	var req *http.Request

	req, e = http.NewRequest(http.MethodPost, url, body)
	if e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// PostForm will make a POST request using WinINet.dll.
func (c *Client) PostForm(
	url string, data url.Values,
) (*http.Response, error) {
	var body io.Reader = bytes.NewReader([]byte(data.Encode()))
	var e error
	var req *http.Request

	req, e = http.NewRequest(http.MethodPost, url, body)
	if e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	return c.Do(req)
}

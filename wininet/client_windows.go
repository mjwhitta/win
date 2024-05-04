//go:build windows

package wininet

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/mjwhitta/errors"
	w32 "github.com/mjwhitta/win/api"
)

// Client is a struct containing relevant metadata to make HTTP
// requests.
type Client struct {
	Debug     bool
	Jar       http.CookieJar
	Timeout   time.Duration
	Transport *http.Transport

	hndl uintptr
	ua   string
}

// NewClient will return a pointer to a new Client instance that
// simply wraps the net/http.Client type.
func NewClient(ua ...string) (*Client, error) {
	var c *Client = &Client{}
	var e error

	if len(ua) == 0 {
		ua = []string{"Go-http-client/1.1"}
	}

	// Store User-Agent
	c.ua = ua[0]

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
	var trans *http.Transport = c.Transport

	if trans == nil {
		if t, ok := http.DefaultTransport.(*http.Transport); ok {
			trans = t
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.ua)
	}

	if reqHndl, e = buildRequest(c.hndl, req, c.Timeout); e != nil {
		return nil, e
	}

	if (trans != nil) && (trans.TLSClientConfig != nil) {
		if trans.TLSClientConfig.InsecureSkipVerify {
			if e = disableTLS(reqHndl); e != nil {
				return nil, e
			}
		}
	}

	if c.Debug {
		// Only needed for debugging b/c WinINet will do this for you
		for _, cookie := range loadCookies(c.Jar, req.URL) {
			req.AddCookie(cookie)
		}

		if b, e = httputil.DumpRequestOut(req, true); e == nil {
			println(string(b))
		}
	}

	if res, e = sendRequest(reqHndl, req); e != nil {
		return nil, e
	}

	if c.Debug {
		if b, e = httputil.DumpResponse(res, true); e == nil {
			println(string(b))
		}
	}

	return res, storeCookies(c.Jar, req.URL)
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

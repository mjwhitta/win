//go:build windows

package winhttp

import (
	"bytes"
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
	c.hndl, e = w32.WinHTTPOpen(
		strings.Join(ua, " "),
		w32.Winhttp.WinhttpAccessTypeAutomaticProxy,
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
	var e error
	var redirect *url.URL
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

	// Load cookies from cookie jar
	loadCookies(c.Jar, req)

	dbgLog(c.Debug, req)

	if res, e = sendRequest(reqHndl, req); e != nil {
		return nil, e
	}

	dbgLog(c.Debug, res)

	// Store cookies into cookie jar
	if e = storeCookies(c.Jar, req.URL, res.Cookies()); e != nil {
		return nil, e
	}

	// Follow redirects
	if redirect, e = res.Location(); e == nil {
		return c.Get(redirect.String())
	}

	return res, nil
}

// Get will make a GET request using WinHTTP.dll.
func (c *Client) Get(url string) (*http.Response, error) {
	var e error
	var req *http.Request

	if req, e = http.NewRequest(http.MethodGet, url, nil); e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	return c.Do(req)
}

// Head will make a HEAD request using WinHTTP.dll.
func (c *Client) Head(url string) (*http.Response, error) {
	var e error
	var req *http.Request

	if req, e = http.NewRequest(http.MethodHead, url, nil); e != nil {
		return nil, errors.Newf("failed to create request: %w", e)
	}

	return c.Do(req)
}

// Post will make a POST request using WinHTTP.dll.
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

// PostForm will make a POST request using WinHTTP.dll.
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

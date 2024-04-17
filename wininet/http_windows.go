//go:build windows

package wininet

import (
	"io"
	"net/http"
	"net/url"
)

// DefaultClient is the default client similar to net/http.
var DefaultClient *Client

// Get will make a GET request using the DefaultClient.
func Get(url string) (*http.Response, error) {
	return DefaultClient.Get(url)
}

// Head will make a HEAD request using the DefaultClient.
func Head(url string) (*http.Response, error) {
	return DefaultClient.Head(url)
}

func init() {
	DefaultClient, _ = NewClient()
}

// Post will make a POST request using the DefaultClient.
func Post(
	url string, contentType string, body io.Reader,
) (*http.Response, error) {
	return DefaultClient.Post(url, contentType, body)
}

// PostForm will make a POST request using the DefaultClient.
func PostForm(url string, data url.Values) (*http.Response, error) {
	return DefaultClient.PostForm(url, data)
}

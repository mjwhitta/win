//go:build !windows
// +build !windows

package http

import "fmt"

// NewClient is only supported on Windows.
func NewClient() (*Client, error) {
	return &Client{}, fmt.Errorf("unsupported OS")
}

// Do is only supported on Windows.
func (c *Client) Do(r *Request) (*Response, error) {
	return nil, fmt.Errorf("unsupported OS")
}

// Get is only supported on Windows.
func (c *Client) Get(url string) (*Response, error) {
	return nil, fmt.Errorf("unsupported OS")
}

// Head is only supported on Windows.
func (c *Client) Head(url string) (*Response, error) {
	return nil, fmt.Errorf("unsupported OS")
}

// Post is only supported on Windows.
func (c *Client) Post(
	url string,
	contentType string,
	body []byte,
) (*Response, error) {
	return nil, fmt.Errorf("unsupported OS")
}

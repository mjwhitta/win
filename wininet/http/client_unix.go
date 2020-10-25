//+build !windows

package http

import "fmt"

// NewClient is only supported on Windows.
func NewClient() (*Client, error) {
	return nil, fmt.Errorf("Unsupported OS")
}

// Get is only supported on Windows.
func (c *Client) Get(
	dst string,
	headers map[string]string,
) (*Response, error) {
	return nil, fmt.Errorf("Unsupported OS")
}

// Post is only supported on Windows.
func (c *Client) Post(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return nil, fmt.Errorf("Unsupported OS")
}

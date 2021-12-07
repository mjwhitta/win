package http

import "time"

// Client is a struct containing relevant metadata to make HTTP
// requests.
type Client struct {
	hndl            uintptr //lint:ignore U1000 Only used for Windows
	Timeout         time.Duration
	TLSClientConfig struct {
		InsecureSkipVerify bool
	}
}

package http

// Client is a struct containing relevant metadata to make HTTP
// requests.
type Client struct {
	hndl            uintptr
	TLSClientConfig struct {
		InsecureSkipVerify bool
	}
}

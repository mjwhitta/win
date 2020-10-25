package http

// DefaultClient is the default client similar to net/http.
var DefaultClient *Client

// Common HTTP methods.
const (
	MethodGet  string = "GET"
	MethodHead string = "HEAD"
	MethodPost string = "POST"
	MethodPut  string = "PUT"
)

// Get will make a GET request using the DefaultClient.
func Get(dst string, headers map[string]string) (*Response, error) {
	return DefaultClient.Get(dst, headers)
}

// Head will make a HEAD request using the DefaultClient.
func Head(dst string, headers map[string]string) (*Response, error) {
	return DefaultClient.Head(dst, headers)
}

func init() {
	var e error

	if DefaultClient, e = NewClient(); e != nil {
		panic(e)
	}
}

// Post will make a POST request using the DefaultClient.
func Post(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return DefaultClient.Post(dst, headers, data)
}

// Put will make a PUT request using the DefaultClient.
func Put(
	dst string,
	headers map[string]string,
	data []byte,
) (*Response, error) {
	return DefaultClient.Put(dst, headers, data)
}

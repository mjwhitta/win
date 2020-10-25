package http

import "io"

// Response is a struct containing common HTTP response data.
type Response struct {
	Body          io.ReadCloser
	ContentLength int64
	Header        map[string][]string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Status        string
	StatusCode    int
}

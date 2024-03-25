//go:build windows

package winhttp

import "net/http"

// Cookie represents an HTTP cookie sent in the Cookie header of an
// HTTP Request.
type Cookie struct {
	http.Cookie
}

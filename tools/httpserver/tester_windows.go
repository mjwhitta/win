//go:build ignore && windows

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	inet "github.com/mjwhitta/win/wininet"
)

var (
	client *inet.Client
	jar    http.CookieJar
	uri    *url.URL
)

func init() {
	var e error
	var host string

	flag.Parse()
	if flag.NArg() == 0 {
		os.Exit(1)
	}

	host = flag.Arg(0)
	if host == "" {
		host = "localhost:8080"
	}

	if uri, e = url.Parse("http://" + host); e != nil {
		panic(e)
	}

	if client, e = inet.NewClient(); e != nil {
		panic(e)
	}

	client.Debug = flag.NArg() > 1

	jar, _ = cookiejar.New(nil)
	client.Jar = jar
}

func main() {
	if e := send(http.MethodGet, "/path"); e != nil {
		panic(e)
	}

	if e := send(http.MethodPost, "/path/to/login"); e != nil {
		panic(e)
	}

	fmt.Println("### Cookies / ###")
	uri, _ = uri.Parse("/")
	for _, cookie := range jar.Cookies(uri) {
		fmt.Println(cookie.String())
	}

	fmt.Println("### Cookies /path ###")
	uri, _ = uri.Parse("/path")
	for _, cookie := range jar.Cookies(uri) {
		fmt.Println(cookie.String())
	}
}

func send(method string, path string) error {
	var e error
	var req *http.Request
	var res *http.Response

	switch method {
	case http.MethodGet:
		req, e = http.NewRequest(method, uri.String()+path, nil)
	case http.MethodPost:
		req, e = http.NewRequest(
			method,
			uri.String()+path,
			bytes.NewBuffer([]byte("user=admin")),
		)
	}

	if e != nil {
		return e
	}

	if res, e = client.Do(req); e != nil {
		return e
	}
	defer res.Body.Close()

	return nil
}

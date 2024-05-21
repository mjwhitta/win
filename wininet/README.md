# WinINet

## Usage

Minimal example:

```
package main

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    "net/http/cookiejar"

    "github.com/mjwhitta/win/wininet"
)

func main() {
    var body io.Reader = bytes.NewReader([]byte("test"))
    var client *wininet.Client
    var dst string = "http://127.0.0.1:8080/asdf"
    var e error
    var req *http.Request
    var res *http.Response

    if t, ok := http.DefaultTransport.(*http.Transport); ok {
        if t.TLSClientConfig == nil {
            t.TLSClientConfig = &tls.Config{}
        }

        t.TLSClientConfig.InsecureSkipVerify = true
    }

    if res, e = wininet.Get(dst); e != nil {
        panic(e)
    }

    if e = output(res); e != nil {
        panic(e)
    }

    req, e = http.NewRequest(http.MethodPost, dst, body)
    if e != nil {
        panic(e)
    }

    req.AddCookie(&http.Cookie{Name: "chocolatechip", Value: "yum"})
    req.AddCookie(&http.Cookie{Name: "oatmealraisin", Value: "gross"})
    req.AddCookie(&http.Cookie{Name: "snickerdoodle", Value: "best"})

    if client, e = wininet.NewClient("My custom UA"); e != nil {
        panic(e)
    }

    if client.Jar, e = cookiejar.New(nil); e != nil {
        panic(e)
    }

    if res, e = client.Do(req); e != nil {
        panic(e)
    }
    defer res.Body.Close()

    if e = output(res); e != nil {
        panic(e)
    }
}

func output(res *http.Response) error {
    var b []byte
    var e error

    if b, e = io.ReadAll(res.Body); e != nil {
        return e
    }

    fmt.Println(res.Status)

    for k := range res.Header {
        fmt.Printf("%s: %s\n", k, res.Header.Get(k))
    }

    for _, cookie := range res.Cookies() {
        fmt.Printf("%s = %s\n", cookie.Name, cookie.Value)
    }

    if len(b) > 0 {
        fmt.Println(string(b))
    }

    fmt.Println()

    return nil
}
```

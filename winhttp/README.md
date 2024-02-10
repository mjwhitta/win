# WinHTTP

## Usage

Minimal example:

```
package main

import (
    "fmt"
    "io"

    http "github.com/mjwhitta/win/winhttp"
)

func main() {
    var b []byte
    var dst = "http://127.0.0.1:8080/asdf"
    var e error
    var headers = map[string]string{
        "User-Agent": "testing, testing, 1, 2, 3...",
    }
    var req *http.Request
    var res *http.Response

    http.DefaultClient.TLSClientConfig.InsecureSkipVerify = true

    if _, e = http.Get(dst); e != nil {
        panic(e)
    }

    req = http.NewRequest(http.MethodPost, dst, []byte("test"))
    req.AddCookie(&http.Cookie{Name: "chocolatechip", Value: "tasty"})
    req.AddCookie(&http.Cookie{Name: "oatmealraisin", Value: "gross"})
    req.AddCookie(&http.Cookie{Name: "snickerdoodle", Value: "yummy"})
    req.Headers = headers

    if res, e = http.DefaultClient.Do(req); e != nil {
        panic(e)
    }

    if res.Body != nil {
        if b, e = io.ReadAll(res.Body); e != nil {
            panic(e)
        }
    }

    fmt.Println(res.Status)
    for k, vs := range res.Header {
        for _, v := range vs {
            fmt.Printf("%s: %s\n", k, v)
        }
    }
    if len(b) > 0 {
        fmt.Println(string(b))
    }

    if len(res.Cookies()) > 0 {
        fmt.Println()
        fmt.Println("# COOKIEJAR")
    }

    for _, cookie := range res.Cookies() {
        fmt.Printf("%s = %s\n", cookie.Name, cookie.Value)
    }
}
```

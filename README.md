# Win

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/win?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/win)
![License](https://img.shields.io/github/license/mjwhitta/win?style=for-the-badge)

## What is this?

This Go modules started as a simple "drop-in" replacement of
`net/http` so that you can use WinHTTP and WinINet on Windows for
better proxy support, with NTLM authentication. Microsoft recommends
[WinINet over WinHTTP] unless you're writing a Windows service.

It has been expanded to include multiple Windows API functions and
constants. There are also some helpers for converting Go/Windows types
and debugging DLLs.

**Note:** This is probably beta quality at best.

[WinINet over WinHTTP]: https://docs.microsoft.com/en-us/windows/win32/wininet/wininet-vs-winhttp

## How to install

Open a terminal and run the following:

```
$ go get --ldflags "-s -w" --trimpath -u github.com/mjwhitta/win
```

## Usage

Minimal example:

```
package main

import (
    "fmt"
    "io"

    // http "github.com/mjwhitta/win/winhttp"
    http "github.com/mjwhitta/win/wininet"
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

## Links

- [Source](https://github.com/mjwhitta/win)

## TODO

- Mirror `net/http` as close as possible
    - CookieJar for the Client
    - etc...
- WinINet
    - FTP client

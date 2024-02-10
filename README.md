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
constants. There are nested modules for converting Go/Windows types,
debugging DLLs, and user identity management.

**Note:** This is probably beta quality at best.

[WinINet over WinHTTP]: https://docs.microsoft.com/en-us/windows/win32/wininet/wininet-vs-winhttp

## How to install

Open a terminal and run the following:

```
$ go get --ldflags "-s -w" --trimpath -u github.com/mjwhitta/win
```

## Usage

See each nested module's README for usage examples.

## Links

- [Source](https://github.com/mjwhitta/win)

## TODO

- Mirror `net/http` as close as possible
    - CookieJar for the Client
    - etc...
- WinINet
    - FTP client
- User
    - Disable privs
    - Enable privs

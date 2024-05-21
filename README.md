# Win

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/win?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/win)
![License](https://img.shields.io/github/license/mjwhitta/win?style=for-the-badge)

## What is this?

This Go module started as a simple "drop-in" replacement of `net/http`
so that you can use WinHTTP and WinINet on Windows for better support
for proxying and NTLM authentication. Microsoft recommends [WinINet
over WinHTTP] unless you're writing a Windows service.

If you want to use a minimal, cross-plaform HTTP client, I recommend
[inet] which uses this module behind the scenes.

This module has been expanded to also include multiple Windows API
functions and constants. There are nested modules for converting
Go/Windows types, debugging DLLs, and user identity management.

[inet]: https://github.com/mjwhitta/inet
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
- [WinHTTP/WinINet equivalents](https://learn.microsoft.com/en-us/windows/win32/winhttp/porting-wininet-applications-to-winhttp#winhttp-equivalents-to-wininet-functions)

## TODO

- WinINet
    - FTP client

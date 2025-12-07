package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/log"
)

var port uint

func init() {
	cli.Align = true
	cli.Banner = filepath.Base(os.Args[0]) + " [OPTIONS]"

	cli.Info("Super simple HTTP listener.")

	cli.Flag(
		&port,
		"p",
		"port",
		8080, //nolint:mnd // Default non-privileged HTTP port
		"Listen on specified port (default: 8080).",
	)
	cli.Parse()
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	var cookie *http.Cookie = &http.Cookie{
		HttpOnly: true,
		Path:     "/path",
		// Secure:   true, // No TLS while testing, so don't enable
	}

	if b, e := httputil.DumpRequest(req, true); e == nil {
		fmt.Println(string(b))
	}

	w.Header().Add("Location", "/path")

	cookie.Name = "chocolatechip"
	cookie.Value = "delicious"
	w.Header().Add("Set-Cookie", cookie.String())

	cookie.Name = "cookiemonster"
	cookie.Value = "hero"
	w.Header().Add("Set-Cookie", cookie.String())

	cookie.Name = "oatmealraisin"
	cookie.Value = "gross"
	w.Header().Add("Set-Cookie", cookie.String())

	cookie.Name = "snickerdoodle"
	cookie.Value = "best"
	w.Header().Add("Set-Cookie", cookie.String())

	cookie.Name = "sugarcookie"
	cookie.Path = "/"
	cookie.Value = "dough"
	w.Header().Add("Set-Cookie", cookie.String())

	w.WriteHeader(http.StatusFound)
}

func main() {
	var addr string
	var e error
	var mux *http.ServeMux
	var server *http.Server

	addr = fmt.Sprintf("0.0.0.0:%d", port)

	mux = http.NewServeMux()
	mux.HandleFunc("/path", rootHandler)
	mux.HandleFunc("/path/to/login", loginHandler)

	server = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, //nolint:mnd // 10 secs
	}

	log.Infof("Listening on %s", addr)

	e = server.ListenAndServe()

	switch e {
	case nil, http.ErrServerClosed:
	default:
		panic(e)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	var cookie *http.Cookie = &http.Cookie{
		HttpOnly: true,
		Path:     "/path",
		// Secure:   true, // No TLS while testing, so don't enable
	}

	if b, e := httputil.DumpRequest(req, true); e == nil {
		fmt.Println(string(b))
	}

	cookie.Name = "chocolatechip"
	if _, e := req.Cookie("chocolatechip"); e != nil {
		cookie.Value = "unknown"
	} else {
		cookie.Value = "yum"
	}

	w.Header().Add("Set-Cookie", cookie.String())

	cookie.Name = "cookiemonster"
	if _, e := req.Cookie("cookiemonster"); e != nil {
		cookie.Value = "unknown"
		w.Header().Add("Set-Cookie", cookie.String())
	}

	cookie.Name = "snickerdoodle"
	if _, e := req.Cookie("snickerdoodle"); e != nil {
		cookie.Value = "unknown"
		w.Header().Add("Set-Cookie", cookie.String())
	}

	_, _ = w.Write([]byte("Success"))
}

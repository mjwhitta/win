package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/log"
)

var port uint

func init() {
	cli.Align = true
	cli.Banner = hl.Sprintf("%s [OPTIONS]", os.Args[0])
	cli.Info("Super simple HTTP listener.")
	cli.Flag(
		&port,
		"p",
		"port",
		8080,
		"Listen on specified port (default: 8080).",
	)
	cli.Parse()
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if b, e := httputil.DumpRequest(req, true); e == nil {
		fmt.Println(string(b))
	}

	w.Header().Add("Location", "/")
	w.Header().Add("Set-Cookie", "chocolatechip=delicious; Path=/")
	w.Header().Add("Set-Cookie", "cookiemonster=hero; Path=/")
	w.Header().Add("Set-Cookie", "oatmealraisin=gross; Path=/")
	w.Header().Add("Set-Cookie", "snickerdoodle=best; Path=/")

	w.WriteHeader(http.StatusFound)
}

func main() {
	var addr string
	var e error
	var mux *http.ServeMux
	var server *http.Server

	addr = hl.Sprintf("0.0.0.0:%d", port)

	mux = http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/some/path/to/login", loginHandler)

	server = &http.Server{Addr: addr, Handler: mux}

	log.Infof("Listening on %s", addr)
	e = server.ListenAndServe()

	switch e {
	case nil, http.ErrServerClosed:
	default:
		panic(e)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	if b, e := httputil.DumpRequest(req, true); e == nil {
		fmt.Println(string(b))
	}

	if _, e := req.Cookie("chocolatechip"); e != nil {
		w.Header().Add("Set-Cookie", "chocolatechip=unknown")
	} else {
		w.Header().Add("Set-Cookie", "chocolatechip=yum")
	}

	if _, e := req.Cookie("cookiemonster"); e != nil {
		w.Header().Add("Set-Cookie", "cookiemonster=unknown")
	}

	if _, e := req.Cookie("snickerdoodle"); e != nil {
		w.Header().Add("Set-Cookie", "snickerdoodle=unknown")
	}

	w.Write([]byte("Success"))
}

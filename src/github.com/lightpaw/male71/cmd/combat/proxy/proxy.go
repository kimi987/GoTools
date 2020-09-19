package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	port      int64
	targetUrl string
)

func main() {
	flag.Int64Var(&port, "p", 8211, "port (default: 8211)")
	flag.StringVar(&targetUrl, "tu", "http://192.168.1.69:7888", "target url (default: http://192.168.1.69:7888)")
	flag.Parse()

	fmt.Println(targetUrl)
	u, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
		return
	}

	http.Handle("/", httputil.NewSingleHostReverseProxy(u))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}

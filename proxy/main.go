package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	serverFirst = "http://localhost:8080"
	serverSecond = "http://localhost:9090"
	count = 0
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if count == 0 {
		url, err := url.Parse(serverFirst)
		if err != nil {
			log.Println(err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
		count++
	} else {
		url, err := url.Parse(serverSecond)
		if err != nil {
			log.Println(err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
		count--
	}
}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}


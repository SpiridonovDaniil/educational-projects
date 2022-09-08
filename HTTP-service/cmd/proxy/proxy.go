package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

const proxyAddr string = ":9001"

var (
	counter = int32(0)
	firstInstanceHost = "http://app1:8080"
	secondInstanceHost = "http://app2:8080"
)

func main() {
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {

	if counter == 0 {
		atomic.AddInt32(&counter, 1)
		url, err := url.Parse(firstInstanceHost)
		if err != nil {
			log.Println(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
		log.Println("1 server processed the request")
		return
	}

	atomic.AddInt32(&counter, -1)
	url, err := url.Parse(secondInstanceHost)
	if err != nil {
		log.Println(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
	log.Println("2 server processed the request")
}
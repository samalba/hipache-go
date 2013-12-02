package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Proxy struct {
	cache *Cache
}

func NewProxy() *Proxy {
	return &Proxy{cache: NewCache()}
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	director := func(target *http.Request) {
		originalUrl := r.URL.String()
		hostHeader := "www.docker.com:80"
		target.URL.Scheme = "http"
		target.URL.Host = hostHeader
		target.Host = hostHeader
		//TODO(samalba): do real http logging
		log.Printf("%s -> %s\n", originalUrl, target.URL)
	}
	p := &httputil.ReverseProxy{Director: director}
	p.ServeHTTP(w, r)
}

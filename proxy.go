package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	cache *Cache
	config *Config
}

func NewProxy(config *Config) *Proxy {
	return &Proxy{NewCache(config), config}
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	director := func(target *http.Request) {
		originalUrl := r.URL.String()
		hostHeader := r.Host
		backend, err := proxy.cache.GetBackend(hostHeader)
		if err != nil {
			log.Println("cache.GetBackend:", err)
			return
		}
		target.URL, err = url.Parse(backend.URL)
		if err != nil {
			//handle wrong URL backend error
			return
		}
		target.URL.Host = hostHeader
		target.Host = hostHeader
		//TODO(samalba): do real http logging
		log.Printf("%s -> %s\n", originalUrl, target.URL)
	}
	p := &httputil.ReverseProxy{Director: director}
	p.ServeHTTP(w, r)
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func runHTTPServer(config *Config) {
	address := fmt.Sprintf("0.0.0.0:%d", config.Server.Port)
	log.Println("Started:", address)
	proxy := NewProxy(config)
	err := http.ListenAndServe(address, proxy)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func main() {
	config := NewConfig()
	config.Parse()
	runHTTPServer(config)
}

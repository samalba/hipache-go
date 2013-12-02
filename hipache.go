package main

import (
	"log"
	"net/http"
)

func runHTTPServer() {
	log.Println("Started: 0.0.0.0:1080")
	proxy := NewProxy()
	err := http.ListenAndServe(":1080", proxy)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func main() {
	//FIXME(samalba): handle config to choose the port (at least)
	runHTTPServer()
}

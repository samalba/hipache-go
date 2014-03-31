package main

import (
	"fmt"
	"log"
	"net/http"
)

func runHTTPServer(config *Config) {
	address := ""
	if config.Server.Port == 0 {
		address = ":80"
	} else {
		address = fmt.Sprintf(":%d", config.Server.Port)
	}
	proxy := NewProxy(config)
	// Enable SSL if needed
	if config.Server.Https.Key != "" && config.Server.Https.Cert != "" {
		go func() {
			if config.Server.Https.Port == 0 {
				address = ":443"
			} else {
				address = fmt.Sprintf(":%d", config.Server.Https.Port)
			}
			log.Println("ListenTLS on", address)
			err := http.ListenAndServeTLS(address, config.Server.Https.Cert,
				config.Server.Https.Key, proxy)
			if err != nil {
				log.Fatal("http.ListenAndServeTLS: ", err)
			}
		}()
	} else {
		log.Println("https is disabled")
	}
	log.Println("Listen on", address)
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

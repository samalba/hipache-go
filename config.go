package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Debug     bool
		AccessLog string
		Port      int
		Https     struct {
			Port int
			Key  string
			Cert string
		}
	}
	RedisHost     string
	RedisPort     int
	RedisDatabase int
	RedisPassword string
}

func NewConfig() *Config {
	c := &Config{}
	c.RedisHost = "localhost"
	c.RedisPort = 6379
	c.RedisDatabase = 0
	return c
}

func (config *Config) parseConfigFile(configFile string) {
	fp, err := os.Open(configFile)
	if err != nil {
		log.Fatal("config.parseConfigFile: ", err)
	}
	defer fp.Close()
	data, _ := ioutil.ReadAll(fp)
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("config.parseConfigFile: ", err)
	}
}

func (config *Config) Parse() {
	configFile := flag.String("config", "", "Config file (json)")
	serverPort := flag.Int("port", 1080, "Server port")
	serverDebug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()
	config.Server.Debug = *serverDebug
	config.Server.Port = *serverPort
	if *configFile != "" {
		log.Println("Reading config file:", *configFile)
		config.parseConfigFile(*configFile)
	}
	log.Printf("Loaded config: %#v\n", config)
}

package main

import (
	"github.com/fzzy/radix/redis"
	"strings"
	"time"
	"log"
)

type Cache struct {
	redisConn *redis.Client
}

type Backend struct {
	// Backend ID (index position from the cache)
	Id int
	// Associated frontend hostname
	Frontend string
	// VritualHost set for this frontend
	VirtualHost string
	// Number of backends in the pool
	PoolSize int
}

func NewCache() *Cache {
	redisConn, err := redis.DialTimeout("tcp", "127.0.0.1:6379",
		time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal("redis.DialTimeout: ", err)
	}
	return &Cache{redisConn: redisConn}
}

func parseHostHeader(hostHeader string) string {
	strings.Parse
}

func (cache *Cache) GetBackend(hostHeader string) (*Backend, error) {
	host = parseHostHeader(hostHeader)
	r := cache.redisConn
	r.Append("MULTI")
	r.Append("LRANGE", "frontend:" + hostHeader, 0, -1)
	r.Append("LRANGE", "frontend:*" + getDomainName(hostHeader), 0, -1)
	r.Append("LRANGE", "frontend:*", 0, -1)
	r.Append("EXEC")
	r.GetReply
}

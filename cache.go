package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Cache struct {
	redisConn *redis.Client
	random    *rand.Rand
	config    *Config
}

type Backend struct {
	// Backend ID (index position from the cache)
	Id int
	// Associated frontend hostname
	Frontend string
	// VritualHost set for this frontend
	VirtualHost string
	// Backend URL to send the request to
	URL string
	// Number of backends in the pool
	PoolSize int
}

func NewCache(config *Config) *Cache {
	address := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	redisConn, err := redis.DialTimeout("tcp", address,
		time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal("redis.DialTimeout: ", err)
	}
	if config.RedisDatabase > 0 {
		r := redisConn.Cmd("select", config.RedisDatabase)
		if r.Err != nil {
			log.Fatal("redis.select: ", r.Err)
		}
	}
	if config.RedisPassword != "" {
		r := redisConn.Cmd("auth", config.RedisPassword)
		if r.Err != nil {
			log.Fatal("redis.auth: ", r.Err)
		}
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Cache{redisConn, random, config}
}

func parseHostHeader(hostHeader string) string {
	parts := strings.SplitN(hostHeader, ":", 2)
	hostHeader = parts[0]
	hostHeader = strings.ToLower(hostHeader)
	return hostHeader
}

// Returns the domain name (without the subdomain)
func getDomainName(hostname string) string {
	parts := strings.SplitAfter(hostname, ".")
	ln := len(parts)
	if ln <= 2 {
		return hostname
	}
	return strings.Join(parts[ln-2:], "")
}

func findReply(replies []*redis.Reply) *redis.Reply {
	for i := 0; i < len(replies); i++ {
		if len(replies[i].Elems) > 0 {
			return replies[i]
		}
	}
	return nil
}

func (cache *Cache) GetBackend(hostHeader string) (*Backend, error) {
	hostHeader = parseHostHeader(hostHeader)
	r := cache.redisConn
	r.Append("multi")
	r.Append("lrange", "frontend:"+hostHeader, 0, -1)
	r.Append("lrange", "frontend:*."+getDomainName(hostHeader), 0, -1)
	r.Append("lrange", "frontend:*", 0, -1)
	r.Append("exec")
	for i := 0; i < 4; i++ {
		// Only the reply of "EXEC" is relevant
		r.GetReply()
	}
	reply := findReply(r.GetReply().Elems)
	if reply == nil {
		return nil, fmt.Errorf("Cannot find a valid backend")
	}
	elems, _ := reply.List()
	// pickup a random backendURL index
	poolSize := len(elems) - 1
	elemIdx := 1
	if poolSize > 1 {
		elemIdx = cache.random.Intn(poolSize) + 1
	}
	backend := &Backend{elemIdx - 1, hostHeader, elems[0], elems[elemIdx],
		poolSize}
	return backend, nil
}

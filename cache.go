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

func NewCache() *Cache {
	redisConn, err := redis.DialTimeout("tcp", "127.0.0.1:6379",
		time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal("redis.DialTimeout: ", err)
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Cache{redisConn: redisConn, random: random}
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
	r.Append("MULTI")
	r.Append("LRANGE", "frontend:"+hostHeader, 0, -1)
	r.Append("LRANGE", "frontend:*."+getDomainName(hostHeader), 0, -1)
	r.Append("LRANGE", "frontend:*", 0, -1)
	r.Append("EXEC")
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

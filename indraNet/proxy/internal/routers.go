package internal

import (
	"cosmo-proxy/internal/rdb"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func GetRouter(addr string) (string, error) {

	cache := rdb.GetClient()

	if cache.RedisClient == nil || cache.MemoryCache == nil {
		log.Panic("Redis client or LRU cache is not initialized")
		return "", fmt.Errorf("Redis client or LRU cache is not initialized")
	}

	// Check if the address is already cached
	cachedAddr, found := cache.MemoryCache.Get(addr)
	if found {
		log.Printf("Cache hit for address: %s", addr)
		fmt.Println("Cache hit for address:", addr)
		return cachedAddr.(string), nil
	}

	redisAddr, err := cache.RedisClient.Get(cache.Ctx, addr).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found in Redis", addr)
	} else if err != nil {
		return "", err
	}

	cache.MemoryCache.Add(addr, redisAddr)

	return redisAddr, nil
}

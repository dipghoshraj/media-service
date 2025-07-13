package rdb

import (
	"context"
	"cosmo-proxy/config"
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/redis/go-redis/v9"
)

var (
	client   *redis.Client
	ctx      = context.Background()
	lruCache *lru.Cache
)

type Cache struct {
	RedisClient *redis.Client
	MemoryCache *lru.Cache
	Ctx         context.Context
}

// Init initializes the Redis client with the provided address and password.

func Init() error {
	client = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"), // Replace with your Redis server address
		Password: config.GetEnv("REDIS_PASSWORD", ""),           // No password set
		DB:       0,                                             // Use default DB
	})

	newCache, err := lru.New(100)
	if err != nil {
		return fmt.Errorf("failed to create LRU cache: %v", err)
	}
	lruCache = newCache

	if err = client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return nil
}

func GetClient() Cache {
	if client == nil && lruCache == nil {
		if err := Init(); err != nil {
			panic(fmt.Sprintf("Failed to initialize Redis client: %v", err))
		}
	}
	return Cache{
		RedisClient: client,
		MemoryCache: lruCache,
		Ctx:         ctx,
	}
}

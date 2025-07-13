package persist

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {
	// Initialize the RDB package
	// This is a placeholder for any initialization logic needed for the RDB package
	rdb = redis.NewClient(&redis.Options{
		Addr:     "3.109.58.10:6379", // Replace with your Redis server address
		Password: "",                 // No password set
		DB:       0,                  // Use default DB
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Dstore: %v", err)
	}
	log.Println("Dstore initialized successfully")
}

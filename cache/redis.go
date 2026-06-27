package cache

import (
	"context"
	"log"
	"time"

	"backend/config"

	"github.com/redis/go-redis/v9"
)

// RDB is the global Redis client used as the durable page cursor for the sync.
var RDB *redis.Client

// Connect opens the Redis connection. A failure here is non-fatal: the server
// keeps running, but the sync workers will error and retry until Redis is up.
func Connect(cfg *config.Config) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Printf("⚠ redis not reachable at %s: %v (sync cursor unavailable until it is)", cfg.RedisAddr, err)
		return
	}
	log.Printf("✓ connected to Redis at %s", cfg.RedisAddr)
}

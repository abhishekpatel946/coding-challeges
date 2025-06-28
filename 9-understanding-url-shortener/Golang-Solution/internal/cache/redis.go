package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Global Redis client
var RedisClient *redis.Client

// InitRedis initializes Redis connection
func InitRedis() error {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 20, // connection pool size
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis not available, using in-memory cache only: %v", err)
		RedisClient = nil
		return nil
	}

	log.Println("✅ Successfully connected to Redis")
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

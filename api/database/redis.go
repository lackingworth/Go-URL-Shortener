package database

import (
	"context"
	
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

// Initializing the Redis Client
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: 		"redis:6379", // For local testing change to local address
		Password: 	"",
		DB: 		dbNo,
	})
	return rdb
}


package database

import (
	"context"
	
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

// Initializing the Redis Client
func CreateClient(dbNo int, addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: 		addr, // For local testing change to local address
		Password: 	"",
		DB: 		dbNo,
	})
	return rdb
}


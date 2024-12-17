package db

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func Connect(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	return redis.NewClient(opt)
}

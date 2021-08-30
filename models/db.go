package models

import "github.com/go-redis/redis"

var client *redis.Client

func InitializeDb() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

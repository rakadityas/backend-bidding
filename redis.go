package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	RedisClient *redis.Client
)

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := RedisClient.Ping().Result()
	fmt.Println(pong, err)
}

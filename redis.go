package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	Client *redis.Client
)

func initRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)

}

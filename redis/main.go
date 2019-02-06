package redis

import (
	"log"

	"github.com/go-redis/redis"
)

// Client Client
var Client *redis.Client

// ConnectClient ConnectClient
func ConnectClient(addr string, password string, DB int) {
	// 建立連線
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB, // use default DB
	})

	_, err := Client.Ping().Result()
	if err != nil {
		log.Fatalf("Ping Redis Error: %v", err)
	}
}

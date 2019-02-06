package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	redis "boardgame_gamecenter/redis"
)

func connectRedis() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")
	port := os.Getenv("REDIS_PORT")
	addr := host + ":" + port

	redis.ConnectClient(addr, password, 0)
}

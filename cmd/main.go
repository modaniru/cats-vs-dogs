package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modaniru/bio-vue/cat-vs-dogs/internal/server"
	"github.com/modaniru/bio-vue/cat-vs-dogs/internal/storage/rstorage"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	storage := rstorage.NewRedisStorage(redisClient)
	server := server.NewServer(storage)
	router := server.InitRouter()

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

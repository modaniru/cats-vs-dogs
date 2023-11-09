package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/modaniru/bio-vue/cat-vs-dogs/internal/server"
	"github.com/modaniru/bio-vue/cat-vs-dogs/internal/storage/rstorage"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
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

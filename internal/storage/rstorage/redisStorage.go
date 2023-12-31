package rstorage

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *redisStorage {
	_, err := client.Get(context.Background(), "cat").Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Fatal(err.Error())
		}
		client.Set(context.Background(), "cat", 0, 0)
	}
	_, err = client.Get(context.Background(), "dog").Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Fatal(err.Error())
		}
		client.Set(context.Background(), "dog", 0, 0)
	}

	return &redisStorage{client: client}
}

func (r *redisStorage) Increase(ctx context.Context, key string) error {
	count, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	count++
	_, err = r.client.Set(ctx, key, count, -1).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisStorage) Get(ctx context.Context, key string) (int, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		return -1, err
	}

	return count, err
}

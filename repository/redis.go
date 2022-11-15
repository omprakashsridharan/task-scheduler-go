package repository

import (
	"context"
	"log"
	"task-scheduler/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(redisConfig config.Redis, ctx context.Context) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Url,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Panicf("Error whole instantiating redis client: %s", err)
	}
	log.Println("Redis client instantiated")
	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}
}

func (rs *RedisStorage) Set(key string, value string, expiration time.Duration) error {
	return rs.client.Set(rs.ctx, key, value, expiration).Err()
}

func (rs *RedisStorage) Get(key string) (string, error) {
	return rs.client.Get(rs.ctx, key).Result()
}

func (rs *RedisStorage) delete(key string) error {
	return rs.client.Del(rs.ctx, key).Err()
}

func (rs *RedisStorage) exists(key string) (bool, error) {
	numberOfFoundKeys, err := rs.client.Exists(rs.ctx, key).Result()
	return numberOfFoundKeys == 1, err
}

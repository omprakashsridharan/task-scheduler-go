package repository

import (
	"context"
	"errors"
	"log"
	"task-scheduler/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

var RedisInstantiationError = errors.New("error whole instantiating redis client")
var RedisParseUrlError = errors.New("error whole parsing redis url")

func NewRedisStorage(redisConfig config.Redis, ctx context.Context) (*RedisStorage, error) {
	options, err := redis.ParseURL(redisConfig.URI)
	if err != nil {
		log.Println(err)
		return nil, RedisParseUrlError
	}
	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, RedisInstantiationError
	}
	log.Println("Redis client instantiated")
	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
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

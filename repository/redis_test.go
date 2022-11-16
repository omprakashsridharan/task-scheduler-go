package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"task-scheduler/config"
	"testing"
	"time"
)

type redisContainer struct {
	testcontainers.Container
	URI string
}

func setupRedis(ctx context.Context) (*redisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())

	return &redisContainer{Container: container, URI: uri}, nil
}

func flushRedis(ctx context.Context, client redis.Client) error {
	return client.FlushAll(ctx).Err()
}

func TestNewRedisStorage(t *testing.T) {
	type args struct {
		redisConfig config.Redis
		ctx         context.Context
	}
	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		t.Log("terminating container")
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: :%s", err)
		}
	})

	options, err := redis.ParseURL(redisContainer.URI)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(options)
	defer func(ctx context.Context, client redis.Client) {
		err := flushRedis(ctx, client)
		if err != nil {

		}
	}(ctx, *client)

	tests := []struct {
		name    string
		args    args
		wantErr func(err error) bool
	}{
		{
			name: "Valid active redis url",
			args: args{
				redisConfig: config.Redis{URI: redisContainer.URI},
				ctx:         ctx,
			},
			wantErr: func(err error) bool {
				return false
			},
		},
		{
			name: "Valid inactive redis url",
			args: args{
				redisConfig: config.Redis{URI: "redis://localhost:1234"},
				ctx:         ctx,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, RedisInstantiationError)
			},
		},
		{
			name: "Invalid redis url",
			args: args{
				redisConfig: config.Redis{URI: "INVALID"},
				ctx:         ctx,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, RedisParseUrlError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedisStorage(tt.args.redisConfig, tt.args.ctx)
			if (err != nil) != tt.wantErr(err) {
				t.Errorf("NewRedisStorage() error = %v, wantErr %v", err, tt.wantErr(err))
				return
			}
		})
	}
}

func TestRedisStorage_Get(t *testing.T) {
	type fields struct {
		client *redis.Client
		ctx    context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
			}
			got, err := rs.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisStorage_Set(t *testing.T) {
	type fields struct {
		client *redis.Client
		ctx    context.Context
	}
	type args struct {
		key        string
		value      string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
			}
			if err := rs.Set(tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisStorage_delete(t *testing.T) {
	type fields struct {
		client *redis.Client
		ctx    context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
			}
			if err := rs.delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisStorage_exists(t *testing.T) {
	type fields struct {
		client *redis.Client
		ctx    context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
			}
			got, err := rs.exists(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

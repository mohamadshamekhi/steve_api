package redisadapter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client *redis.Client
	cfg    Config
}

func New(config Config) *Adapter {
	return &Adapter{client: redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})}
}

func (a *Adapter) Client() *redis.Client {
	return a.client
}

func (a *Adapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}

func (a *Adapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return a.client.Set(ctx, key, value, expiration).Err()
}

func (a *Adapter) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return a.client.Del(ctx, keys...).Err()
}

package cache

import (
	"context"
	"time"
)

type Provider interface {
	Get(ctx context.Context, key string) (string, error)

	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	Delete(ctx context.Context, keys ...string) error
}

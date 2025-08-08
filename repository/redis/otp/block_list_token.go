package redisotp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (d DB) AddToBlocklist(ctx context.Context, kind string, JWTID string, ExpireAt time.Time) error {
	ttl := time.Until(ExpireAt)
	return d.adapter.Set(ctx, fmt.Sprintf("%s:%s", kind, JWTID), "revoked", ttl)
}

func (d DB) IsInBlocklist(ctx context.Context, kind string, JWTID string) (bool, error) {
	_, err := d.adapter.Get(ctx, fmt.Sprintf("%s:%s", kind, JWTID))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return false, err
}

package redisotp

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/redis/go-redis/v9"
	"time"
)

func (d DB) CheckUserHasActiveCode(ctx context.Context, key string) (bool, error) {
	const op = "redisotp.CheckUserHasActiveCode"
	isExist, err := d.adapter.Client().Exists(ctx, key).Result()
	if err != nil {
		return false, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithMeta(map[string]any{"key": key})
	}

	if isExist == 0 {
		return false, nil
	}
	return true, nil
}

func (d DB) GetUserActiveCode(ctx context.Context, key string) (string, error) {
	const op = "redisotp.GetUserActiveCode"
	res, err := d.adapter.Client().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnauthorized).WithMessageKey("validation.incorrect_opt_code")
		}
		return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithMeta(map[string]any{"key": key})
	}
	return res, nil
}

func (d DB) StoreNewActiveCodeForUser(ctx context.Context, key string, code string, expireTime time.Duration) error {
	const op = "redisotp.StoreNewActiveCodeForUser"

	err := d.adapter.Client().Set(ctx, key, code, expireTime).Err()
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithMeta(map[string]any{"key": key, "code": code})
	}
	return nil
}

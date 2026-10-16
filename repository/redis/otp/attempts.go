package redisotp

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// IncrAttempts bumps a counter under key, starting its ttl on the first hit,
// and returns the new count. Used to cap OTP guesses.
func (d DB) IncrAttempts(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	const op = "redisotp.IncrAttempts"

	n, err := d.adapter.Client().Incr(ctx, key).Result()
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	if n == 1 {
		d.adapter.Client().Expire(ctx, key, ttl)
	}
	return n, nil
}

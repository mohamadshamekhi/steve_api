package mongodbsell

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Add(ctx context.Context, sell entity.Sell) error {
	const op = "mongodbsell.Create"
	_, err := d.collection.InsertOne(ctx, sell)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"sell": sell}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

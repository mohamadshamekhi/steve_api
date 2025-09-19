package mongodborder

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, order entity.Order) error {
	const op = "mongodborder.Create"
	_, err := d.collection.InsertOne(ctx, order)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"order": order}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

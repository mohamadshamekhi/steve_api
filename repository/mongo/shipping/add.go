package mongodbshipping

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, sh entity.Shipping) error {
	const op = "mongodbshipping.Create"
	_, err := d.collection.InsertOne(ctx, sh)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"shipping": sh}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

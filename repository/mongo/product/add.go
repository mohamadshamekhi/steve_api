package mongodbproduct

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, product entity.Product) error {
	const op = "mongodbproduct.Create"
	_, err := d.collection.InsertOne(ctx, product)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]any{"product": product}).WithKind(richerror.KindUnexpected)
	}

	return nil
}

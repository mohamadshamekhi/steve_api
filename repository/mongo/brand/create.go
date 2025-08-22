package mongodbbrand

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, brand entity.Brand) error {
	const op = "mongodbbrand.Create"
	_, err := d.collection.InsertOne(ctx, brand)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"brand": brand}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

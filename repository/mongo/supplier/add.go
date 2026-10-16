package mongodbsupplier

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, s entity.Supplier) error {
	const op = "mongodbsupplier.Create"

	if _, err := d.collection.InsertOne(ctx, s); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

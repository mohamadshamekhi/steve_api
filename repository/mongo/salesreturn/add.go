package mongodbsalesreturn

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, r entity.SalesReturn) error {
	const op = "mongodbsalesreturn.Create"

	if _, err := d.collection.InsertOne(ctx, r); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

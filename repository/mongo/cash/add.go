package mongodbcash

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, t entity.CashTransaction) error {
	const op = "mongodbcash.Create"

	if _, err := d.collection.InsertOne(ctx, t); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

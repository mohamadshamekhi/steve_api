package mongodbpurchase

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, p entity.Purchase) error {
	const op = "mongodbpurchase.Create"

	if _, err := d.collection.InsertOne(ctx, p); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

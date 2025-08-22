package mongodbdiscount

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, discount entity.Discount) error {
	const op = "mongodbdiscount.Create"
	_, err := d.collection.InsertOne(ctx, discount)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"discount": discount}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

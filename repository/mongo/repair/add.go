package mongodbrepair

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, repair entity.Repair) error {
	const op = "mongodbrepair.Create"
	_, err := d.collection.InsertOne(ctx, repair)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"brand": repair}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

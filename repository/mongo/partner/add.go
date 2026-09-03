package mongodbpartner

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, p entity.Partner) error {
	const op = "mongodbpartner.Create"
	if _, err := d.collection.InsertOne(ctx, p); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"partner": p})
	}
	return nil
}

package mongodbcomment

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, cm entity.Comment) error {
	const op = "mongodbcomment.Create"
	_, err := d.collection.InsertOne(ctx, cm)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"cm": cm}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

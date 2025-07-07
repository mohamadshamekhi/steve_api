package mongodbcategory

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, category entity.Category) error {
	const op = "mongodbcategory.Create"
	_, err := d.collection.InsertOne(ctx, category)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]any{"category": category}).WithKind(richerror.KindUnexpected)
	}

	return nil
}

package mongodbrole

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, role *entity.Role) error {
	const op = "mongodbrole.Create"
	_, err := d.collection.InsertOne(ctx, role)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"role": role}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

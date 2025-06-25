package mongodbpermission

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, permission *entity.Permission) error {
	const op = "mongodbpermission.Create"
	_, err := d.collection.InsertOne(ctx, permission)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"permission": permission}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

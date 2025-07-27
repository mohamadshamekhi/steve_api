package mongodbuser

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) AddNewUser(ctx context.Context, u entity.User) error {
	const op = "mongodbuser.AddNewUser"
	_, err := d.collection.InsertOne(ctx, u)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"user": u}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

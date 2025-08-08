package mongodbpermission

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) FindByName(ctx context.Context, name string) (entity.Permission, error) {
	const op = "mongodbpermission.FindByName"
	filter := bson.M{"name": name}
	permission := entity.Permission{}
	if err := d.collection.FindOne(ctx, filter).Decode(&permission); err != nil {
		return permission, richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}
	return permission, nil
}

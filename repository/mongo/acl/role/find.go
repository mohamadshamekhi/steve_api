package mongodbrole

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) FindByName(ctx context.Context, name string) (entity.Role, error) {
	const op = "mongodbrole.FindByName"
	filter := bson.M{"name": name}
	var role entity.Role
	if err := d.collection.FindOne(ctx, filter).Decode(&role); err != nil {
		return role, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	return role, nil
}

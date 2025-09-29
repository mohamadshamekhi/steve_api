package mongodbrole

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) GetAllRoles(ctx context.Context, search string) ([]entity.Role, error) {
	const op = "mongodbpermission.GetAllRoles"

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	if search != "" {
		filter["description"] = bson.M{"$regex": search, "$options": "i"}
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var roles []entity.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return roles, nil

}

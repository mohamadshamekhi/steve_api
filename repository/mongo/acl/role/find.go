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

func (d DB) GetPermissionIDsOfRole(ctx context.Context, roleIDs []string) ([]string, error) {
	const op = "mongodbrole.GetPermissionIDsOfRole"
	if len(roleIDs) == 0 {
		return []string{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": roleIDs}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	permissionIDSet := make(map[string]struct{})

	for cursor.Next(ctx) {
		var role entity.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, richerror.New(op).WithWarpError(err).
				WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
		}

		// ۴. اضافه کردن تمام شناسه‌های دسترسی این نقش به map
		for _, permID := range role.PermissionIDs {
			permissionIDSet[permID] = struct{}{}
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	result := make([]string, 0, len(permissionIDSet))
	for id := range permissionIDSet {
		result = append(result, id)
	}

	return result, nil
}

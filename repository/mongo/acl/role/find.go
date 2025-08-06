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

func (d DB) GetPermissionKeysOfRole(ctx context.Context, roleKeys []string) ([]string, error) {
	const op = "mongodbrole.GetPermissionIDsOfRole"
	if len(roleKeys) == 0 {
		return []string{}, nil
	}
	filter := bson.M{"name": bson.M{"$in": roleKeys}}
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	permissionKsySet := make(map[string]struct{})

	for cursor.Next(ctx) {
		var role entity.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, richerror.New(op).WithWarpError(err).
				WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
		}

		for _, permID := range role.PermissionKeys {
			permissionKsySet[permID] = struct{}{}
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	result := make([]string, 0, len(permissionKsySet))
	for id := range permissionKsySet {
		result = append(result, id)
	}
	return result, nil
}

package mongodbrole

import (
	"context"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdatePermissionKeys overwrites the permission_keys array of the role with
// the given name. Used by the seeder to sync newly predefined permissions
// into roles that already exist in the database (FindByName + Create only
// handles roles that don't exist yet, see internal/seeder/acl.go).
func (d DB) UpdatePermissionKeys(ctx context.Context, name string, permissionKeys []string) error {
	const op = "mongodbrole.UpdatePermissionKeys"
	filter := bson.M{"name": name}
	update := bson.M{"$set": bson.M{"permission_keys": permissionKeys}}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}
	return nil
}
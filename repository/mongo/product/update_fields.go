package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateFields writes only the named bson fields of p, so a concurrent order
// can't have its stock change overwritten the way a full ReplaceOne would.
// A field that omitempty drops from p is unset, matching what a replace does.
func (d DB) UpdateFields(ctx context.Context, p entity.Product, fields ...string) error {
	const op = "mongodbproduct.UpdateFields"

	p.UpdatedAt = time.Now()
	raw, err := bson.Marshal(p)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	var doc bson.M
	if err := bson.Unmarshal(raw, &doc); err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	set := bson.M{"updated_at": doc["updated_at"]}
	unset := bson.M{}
	for _, field := range fields {
		if value, ok := doc[field]; ok {
			set[field] = value
		} else {
			unset[field] = ""
		}
	}
	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	filter := bson.M{"_id": p.ID, "deleted_at": bson.M{"$exists": false}}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

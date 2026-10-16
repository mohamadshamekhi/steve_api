package mongodbarticle

import (
	"context"

	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// IncrementViewCount atomically bumps an article's view_count by one.
// Callers are expected to have already deduplicated the visitor (see
// repository/redis/articleview) before calling this.
func (d DB) IncrementViewCount(ctx context.Context, id string) error {
	const op = "mongodbarticle.IncrementViewCount"

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	if _, err := d.collection.UpdateOne(ctx, filter, update); err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithMeta(map[string]any{"id": id})
	}

	return nil
}

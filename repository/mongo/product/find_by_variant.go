package mongodbproduct

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByVariantID(ctx context.Context, variantID string) (entity.Product, error) {
	const op = "mongodbproduct.FindByVariantID"
	filter := bson.M{
		"variants.id": variantID,
		"deleted_at":  bson.M{"$exists": false},
	}
	var p entity.Product

	if err := d.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return p, richerror.New(op).WithWarpError(err).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"variant_id": variantID})
	}

	return p, nil
}

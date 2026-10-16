package mongodbpurchase

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Purchase, error) {
	const op = "mongodbpurchase.FindByID"

	var p entity.Purchase
	if err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Purchase{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("purchase.not_found")
		}
		return entity.Purchase{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return p, nil
}

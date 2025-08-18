package mongodbcart

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Cart, error) {
	const op = "mongodbcart.FindByID"
	filter := bson.M{"_id": id}
	var cart entity.Cart

	if err := d.collection.FindOne(ctx, filter).Decode(&cart); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return cart, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return cart, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return cart, nil
}

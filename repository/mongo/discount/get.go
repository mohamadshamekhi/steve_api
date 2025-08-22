package mongodbdiscount

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Discount, error) {
	const op = "mongodbdiscount.FindByID"
	filter := bson.M{"_id": id}
	var discount entity.Discount

	if err := d.collection.FindOne(ctx, filter).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return discount, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return discount, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return discount, nil
}

func (d DB) CheckCodeIsUnique(ctx context.Context, code string) (bool, error) {
	const op = "mongodbdiscount.CheckCodeIsUnique"

	filter := bson.M{"code": code}
	var discount entity.Discount

	if err := d.collection.FindOne(ctx, filter).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"code": code})
	}
	return false, nil
}

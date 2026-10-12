package mongodbshipping

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByName(ctx context.Context, name string) (entity.Shipping, error) {
	const op = "mongodbshipping.FindByName"
	filter := bson.M{"name": name}
	var sh entity.Shipping

	if err := d.collection.FindOne(ctx, filter).Decode(&sh); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return sh, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return sh, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"name": name})
	}
	return sh, nil
}

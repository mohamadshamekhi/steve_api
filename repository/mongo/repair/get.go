package mongodbrepair

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Repair, error) {
	const op = "mongodbrepair.FindByID"
	filter := bson.M{"_id": id}
	var r entity.Repair
	if err := d.collection.FindOne(ctx, filter).Decode(&r); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return r, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return r, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return r, nil
}

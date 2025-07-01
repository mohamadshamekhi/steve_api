package mongodbmedia

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) CheckMediaIsExist(ctx context.Context, id string) (bool, error) {
	const op = "mongodbmedia.CheckMediaIsExist"
	filter := bson.M{"_id": id}
	var media entity.Media

	if err := d.collection.FindOne(ctx, filter).Decode(&media); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return true, nil
}

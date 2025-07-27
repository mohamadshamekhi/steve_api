package mongodbuser

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) CheckUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	const op = "mongodbuser.IsBUserIsExistByPhoneNumber"

	var user entity.User
	if err := d.collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"phoneNumber": phoneNumber})
	}

	return true, nil
}

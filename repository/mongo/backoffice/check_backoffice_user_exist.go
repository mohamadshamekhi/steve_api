package mongodbbackoffice

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) IsBackOfficeUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	const op = "mongodbbackoffice.IsBackOfficeUserIsExistByPhoneNumber"

	var backOfficeUser entity.BackOffice
	if err := d.collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&backOfficeUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"phoneNumber": phoneNumber})
	}
	if backOfficeUser.ID == "" {
		return false, nil
	}
	return true, nil
}

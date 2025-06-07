package mongodbbackoffice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) IsBackOfficeUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	const op = "mongodbbackoffice.IsBackOfficeUserIsExistByPhoneNumber"

	var backOfficeUser entity.BackOffice
	if err := d.collection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&backOfficeUser); err != nil {
		return false, richerror.New(op).WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"phoneNumber": phoneNumber})
	}
	if backOfficeUser.ID == "" {
		return false, nil
	}
	return true, nil
}

package mongodbbackoffice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) GetBackOfficeUserIDByPhoneNumber(ctx context.Context, phoneNumber string) (string, error) {
	const op = "mongodbbackoffice.GetBackOfficeUserIDByPhoneNumber"
	filter := bson.M{"phone_number": phoneNumber}
	var backOffice entity.BackOffice
	if dErr := d.collection.FindOne(ctx, filter).Decode(&backOffice); dErr != nil {
		return "", richerror.New(op).WithWarpError(dErr).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).
			WithKind(richerror.KindUnexpected)
	}
	return backOffice.ID, nil
}

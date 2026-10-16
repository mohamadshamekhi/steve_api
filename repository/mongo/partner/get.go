package mongodbpartner

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Partner, error) {
	const op = "mongodbpartner.FindByID"
	var p entity.Partner
	if err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("partner.not_found")
		}
		return p, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	return p, nil
}

func (d DB) FindByBackofficeUserID(ctx context.Context, backofficeUserID string) (entity.Partner, error) {
	const op = "mongodbpartner.FindByBackofficeUserID"
	var p entity.Partner
	if err := d.collection.FindOne(ctx, bson.M{"backoffice_user_id": backofficeUserID}).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("partner.not_found")
		}
		return p, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	return p, nil
}

func (d DB) FindByUserID(ctx context.Context, userID string) (entity.Partner, error) {
	const op = "mongodbpartner.FindByUserID"
	var p entity.Partner
	if err := d.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindNotFound).WithMessageKey("partner.not_found")
		}
		return p, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	return p, nil
}

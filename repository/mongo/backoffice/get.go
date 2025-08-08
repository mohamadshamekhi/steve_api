package mongodbbackoffice

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) GetBackOfficeUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.BackOffice, error) {
	const op = "mongodbbackoffice.GetBackOfficeUserByPhonenumber"
	filter := bson.M{"phone_number": phoneNumber}
	var backOffice entity.BackOffice
	if dErr := d.collection.FindOne(ctx, filter).Decode(&backOffice); dErr != nil {
		if errors.Is(dErr, mongo.ErrNoDocuments) {
			return backOffice, richerror.New(op).WithWarpError(dErr).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return backOffice, richerror.New(op).WithWarpError(dErr).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	return backOffice, nil
}

func (d DB) GetBackOfficeUserByID(ctx context.Context, id string) (entity.BackOffice, error) {
	const op = "mongodbbackoffice.GetBackOfficeUserByID"
	filter := bson.M{"_id": id}
	var backOffice entity.BackOffice
	if dErr := d.collection.FindOne(ctx, filter).Decode(&backOffice); dErr != nil {
		if errors.Is(dErr, mongo.ErrNoDocuments) {
			return backOffice, richerror.New(op).WithWarpError(dErr).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return backOffice, richerror.New(op).WithWarpError(dErr).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	return backOffice, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.BackOffice, error) {
	const op = "mongodbbackoffice.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.BackOffice{}, nil
	}


	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var users []entity.BackOffice
	if err := cursor.All(ctx, &users); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return users, nil
}
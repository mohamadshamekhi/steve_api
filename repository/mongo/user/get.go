package mongodbuser

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error) {
	const op = "mongodbuser.GetUserByPhoneNumber"
	filter := bson.M{"phone_number": phoneNumber}
	var user entity.User
	if dErr := d.collection.FindOne(ctx, filter).Decode(&user); dErr != nil {
		if errors.Is(dErr, mongo.ErrNoDocuments) {
			return user, richerror.New(op).WithWarpError(dErr).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return user, richerror.New(op).WithWarpError(dErr).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d DB) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	const op = "mongodbuser.GetUserByID"
	filter := bson.M{"_id": id}
	var user entity.User
	if dErr := d.collection.FindOne(ctx, filter).Decode(&user); dErr != nil {
		if errors.Is(dErr, mongo.ErrNoDocuments) {
			return user, richerror.New(op).WithWarpError(dErr).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return user, richerror.New(op).WithWarpError(dErr).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d DB) FindManyByIDs(ctx context.Context, ids []string) ([]entity.User, error) {
	const op = "mongodbuser.FindManyByIDs"
	if len(ids) == 0 {
		return []entity.User{}, nil
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

	var users []entity.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return users, nil
}

func (d DB) FindNonExistingIDs(ctx context.Context, ids []string) ([]string, error) {
	const op = "mongodbuser.FindNonExistingIDs"

	if len(ids) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	existingIDSet := make(map[string]struct{}, len(results))
	for _, result := range results {
		existingIDSet[result.ID] = struct{}{}
	}

	var notFoundIDs []string
	for _, id := range ids {
		if _, found := existingIDSet[id]; !found {
			notFoundIDs = append(notFoundIDs, id)
		}
	}

	return notFoundIDs, nil

}

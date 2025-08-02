package mongodbuser

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) AddAddress(ctx context.Context, userID string, address entity.Address) error {
	const op = "mongodbuser.AddAddress"

	filter := bson.M{"_id": userID}

	update := bson.M{"$push": bson.M{"addresses": address}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("user.not-found").
			WithMeta(map[string]any{"user_id": userID})
	}

	return nil
}

func (d DB) UpdateAddress(ctx context.Context, userID string, address entity.Address) error {
	const op = "mongodbuser.UpdateAddress"

	filter := bson.M{
		"_id":          userID,
		"addresses.id": address.ID,
	}

	update := bson.M{
		"$set": bson.M{"addresses.$": address},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found").
			WithMeta(map[string]any{"user_id": userID, "address_id": address.ID})
	}

	return nil
}

func (d DB) FindAddress(ctx context.Context, userID string, addressID string) (entity.Address, error) {
	const op = "mongodbuser.FindAddress"
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{"_id": userID}}},
		bson.D{{"$unwind", "$addresses"}},
		bson.D{{"$match", bson.M{"addresses.id": addressID}}},
		bson.D{{"$replaceRoot", bson.M{"newRoot": "$addresses"}}},
	}

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entity.Address{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var address entity.Address
		if err := cursor.Decode(&address); err != nil {
			return entity.Address{}, richerror.New(op).WithWarpError(err).
				WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
		}
		return address, nil
	}

	if err := cursor.Err(); err != nil {
		return entity.Address{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return entity.Address{}, richerror.New(op).
		WithKind(richerror.KindNotFound).
		WithMessageKey("global.not-found").
		WithMeta(map[string]any{"user_id": userID, "address_id": addressID})
}

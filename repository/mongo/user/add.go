package mongodbuser

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) AddNewUser(ctx context.Context, u entity.User) error {
	const op = "mongodbuser.AddNewUser"
	_, err := d.collection.InsertOne(ctx, u)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"user": u}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

func (d DB) GetOrCreateUser(ctx context.Context, phoneNumber string) (entity.User, error) {
	return d.GetOrCreateUserWithProfile(ctx, phoneNumber, "", nil)
}

// GetOrCreateUserWithProfile upserts by phone number; name and birthDate are
// written only when the document is inserted, never over an existing user.
func (d DB) GetOrCreateUserWithProfile(ctx context.Context, phoneNumber, name string, birthDate *time.Time) (entity.User, error) {
	const op = "mongodbuser.GetOrCreateUserWithProfile"

	filter := bson.M{"phone_number": phoneNumber}

	onInsert := bson.M{
		"_id":          ulid.Make().String(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
		"phone_number": phoneNumber,
	}
	if name != "" {
		onInsert["name"] = name
	}
	if birthDate != nil {
		onInsert["birth_date"] = *birthDate
	}
	update := bson.M{"$setOnInsert": onInsert}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	result := d.collection.FindOneAndUpdate(ctx, filter, update, opts)
	if err := result.Err(); err != nil {
		return entity.User{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	var user entity.User
	if err := result.Decode(&user); err != nil {
		return entity.User{}, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return user, nil
}

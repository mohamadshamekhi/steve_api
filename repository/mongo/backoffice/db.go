package mongodbbackoffice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "backofficeUsers"

type DB struct {
	collection *mongo.Collection
}

// todo - fix this !!!
func (d DB) AddRoleToBackofficeUser(phoneNumber string, role string) error {
	filter := bson.M{"phone_number": phoneNumber}
	var backOffice entity.BackOffice
	if dErr := d.collection.FindOne(context.Background(), filter).Decode(&backOffice); dErr != nil {
		return richerror.New("mongodbbackoffice.AddRoleToBackofficeUser").WithWarpError(dErr).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}
	for _, k := range backOffice.Roles {
		if k == role {
			return nil
		}
	}
	backOffice.Roles = append(backOffice.Roles, role)
	update := bson.M{
		"$set": bson.M{"roles": backOffice.Roles},
	}

	_, uErr := d.collection.UpdateOne(context.Background(), filter, update)
	if uErr != nil {
		return richerror.New("mongodbbackoffice.AddRoleToBackofficeUser").WithWarpError(uErr).
			WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected)
	}

	return nil

}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "phone_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

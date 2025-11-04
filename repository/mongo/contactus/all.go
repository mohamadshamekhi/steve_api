package mongodbcontactus

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	contactusservice "github.com/abolfazlnorzad/steve/service/contactus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) All(ctx context.Context, p commonparam.PaginationRequest) (contactusservice.ListResult, error) {
	const op = "mongodbcontactus.All"
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return contactusservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return contactusservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)
	var contacts []entity.ContactUs
	if err := cursor.All(ctx, &contacts); err != nil {
		return contactusservice.ListResult{}, richerror.New(op).WithWarpError(err).
			WithMessageKey("global.something_wrong").WithKind(richerror.KindUnexpected)
	}

	return contactusservice.ListResult{
		Contacts:   contacts,
		TotalCount: totalCount,
	}, nil
}

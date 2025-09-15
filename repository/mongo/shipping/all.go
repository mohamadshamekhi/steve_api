package mongodbshipping

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	settingservice "github.com/abolfazlnorzad/steve/service/setting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetAll(ctx context.Context, f settingparam.FilterShipping, p commonparam.PaginationRequest) (settingservice.ListResult, error) {
	const op = "mongodbshipping.GetAll"

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	if f.IsActive != nil {
		filter["is_active"] = *f.IsActive
	}
	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return settingservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return settingservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var shipping []entity.Shipping
	if err := cursor.All(ctx, &shipping); err != nil {
		return settingservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return settingservice.ListResult{
		Shipping:   shipping,
		TotalCount: totalCount,
	}, nil
}

package mongodbsupplier

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All lists non-deleted suppliers. search matches name/shop name/phone,
// isActive filters by status when not nil.
func (d DB) All(ctx context.Context, search string, isActive *bool, p commonparam.PaginationRequest) ([]entity.Supplier, int64, error) {
	const op = "mongodbsupplier.All"

	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"shop_name": bson.M{"$regex": search, "$options": "i"}},
			{"phone_number": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if isActive != nil {
		filter["is_active"] = *isActive
	}

	total, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	findOptions := options.Find().
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize())).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	suppliers := make([]entity.Supplier, 0)
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, 0, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return suppliers, total, nil
}

package mongodbdiscount

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d DB) FindByID(ctx context.Context, id string) (entity.Discount, error) {
	const op = "mongodbdiscount.FindByID"
	filter := bson.M{"_id": id}
	var discount entity.Discount

	if err := d.collection.FindOne(ctx, filter).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return discount, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return discount, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": id})
	}
	return discount, nil
}

func (d DB) FindByCode(ctx context.Context, code string) (entity.Discount, error) {
	const op = "mongodbdiscount.FindByCode"
	filter := bson.M{"code": code}
	var discount entity.Discount

	if err := d.collection.FindOne(ctx, filter).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return discount, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return discount, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"id": code})
	}
	return discount, nil
}

func (d DB) CheckCodeIsUnique(ctx context.Context, code string) (bool, error) {
	const op = "mongodbdiscount.CheckCodeIsUnique"

	filter := bson.M{"code": code}
	var discount entity.Discount

	if err := d.collection.FindOne(ctx, filter).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, richerror.New(op).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]any{"code": code})
	}
	return false, nil
}

func (d DB) FindApplicableDiscounts(ctx context.Context, code string) ([]entity.Discount, error) {
	const op = "mongodbdiscount.FindApplicableDiscounts"
	now := time.Now()

	codeOrClause := []bson.M{
		{"code": bson.M{"$in": []interface{}{nil, ""}}},
	}

	if code != "" {
		codeOrClause = append(codeOrClause, bson.M{"code": code})
	}

	filter := bson.M{
		"is_active":  true,
		"starts_at":  bson.M{"$lte": now},
		"expires_at": bson.M{"$gte": now},
		"$or":        codeOrClause,
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var discounts []entity.Discount
	if err := cursor.All(ctx, &discounts); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return discounts, nil
}

func (d DB) GetActiveSpecialDiscount(ctx context.Context) (entity.Discount, error) {
	const op = "mongodbdiscount.GetActiveSpecialDiscount"
	var result entity.Discount

	now := time.Now()

	filter := bson.M{
		"is_special": true,
		"is_active":  true,
		"starts_at":  bson.M{"$lte": now},
		"expires_at": bson.M{"$gte": now},
	}

	err := d.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Discount{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
		}
		return entity.Discount{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	return result, nil
}

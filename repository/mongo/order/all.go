package mongodborder

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) GetAll(ctx context.Context, f orderparam.BackofficeFilterOrder, p commonparam.PaginationRequest) (orderservice.ListResult, error) {
	const op = "mongodborder.GetAll"
	filter, empty, err := d.backofficeOrderFilter(ctx, f)
	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if empty {
		return orderservice.ListResult{Orders: []entity.Order{}}, nil
	}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := d.collection.Find(ctx, filter, findOptions)

	if err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var orders []entity.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return orderservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return orderservice.ListResult{
		Orders:     orders,
		TotalCount: totalCount,
	}, nil
}

func (d DB) GetUserOrderWithDetails(ctx context.Context, userID string, p commonparam.PaginationRequest) (orderservice.StorefrontListResult, error) {
	const op = "mongodborder.GetUserOrderWithDetails"

	filter := bson.M{"user_id": userID}

	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return orderservice.StorefrontListResult{}, richerror.New(op).WithWarpError(err)
	}

	if totalCount == 0 {
		return orderservice.StorefrontListResult{Orders: []orderparam.StoreFrontOrderInfo{}, TotalCount: 0}, nil
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // جدیدترین‌ها اول
		SetSkip(int64(p.GetOffset())).
		SetLimit(int64(p.GetPageSize()))

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return orderservice.StorefrontListResult{}, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var dbOrders []entity.Order
	if err := cursor.All(ctx, &dbOrders); err != nil {
		return orderservice.StorefrontListResult{}, richerror.New(op).WithWarpError(err)
	}

	ordersInfo := make([]orderparam.StoreFrontOrderInfo, 0, len(dbOrders))
	for _, o := range dbOrders {
		ordersInfo = append(ordersInfo, orderparam.StoreFrontOrderInfo{
			ID:              o.ID,
			UserID:          o.UserID,
			Items:           o.Items,
			ShippingAddress: o.ShippingAddress,
			SubtotalAmount:  o.SubtotalAmount,
			ShippingAmount:  o.ShippingAmount,
			FinalAmount:     o.FinalAmount,
			Shipping:        o.Shipping,
			Status:          o.Status,
			CreatedAt:       o.CreatedAt,
			UpdatedAt:       o.UpdatedAt,
		})
	}

	return orderservice.StorefrontListResult{
		Orders:     ordersInfo,
		TotalCount: totalCount,
	}, nil
}

func (d DB) GetUserOrderStatics(ctx context.Context, userID string, status entity.OrderStatus) (uint64, error) {
	const op = "mongodborder.GetUserOrderStatics"

	filter := bson.M{
		"user_id":    userID,
		"status":     status,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return uint64(count), nil
}

func (d DB) GetExpiredPendingOrders(ctx context.Context, expiredBefore time.Time) ([]entity.Order, error) {
	const op = "mongodborder.GetExpiredPendingOrders"

	filter := bson.M{
		"status":     entity.OrderStatusPendingPayment,
		"created_at": bson.M{"$lte": expiredBefore},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var orders []entity.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return orders, nil
}

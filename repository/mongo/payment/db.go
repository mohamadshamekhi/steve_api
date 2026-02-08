package mongodbpayment

import (
	"context"
	"errors"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "payments"

type DB struct {
	collection *mongo.Collection
}

func New(c *mongodb.DB) *DB {
	collection := c.Connection().Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "tracking_code", Value: 1}, {Key: "order_id", Value: 1}, {Key: "gateway_transaction_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}
	return &DB{
		collection: collection,
	}
}

func (d DB) CreatePayment(ctx context.Context, p entity.Payment) error {
	const op = "mongodbpayment.CreatePayment"
	_, err := d.collection.InsertOne(ctx, p)
	if err != nil {
		return richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	return nil
}

func (d DB) UpdateGatewayTransactionID(ctx context.Context, paymentID string, transactionID string) error {
	const op = "mongodbpayment.UpdateGatewayTransactionID"

	filter := bson.M{"_id": paymentID}

	update := bson.M{
		"$set": bson.M{
			"gateway_transaction_id": transactionID,
			"updated_at":             time.Now(),
		},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("payment.not_found").
			WithMeta(map[string]any{"payment_id": paymentID})
	}

	return nil
}
func (d DB) GetPaymentByGatewayTransactionID(ctx context.Context, authority string) (entity.Payment, error) {
	const op = "mongodbpayment.GetPaymentByGatewayTransactionID"

	filter := bson.M{"gateway_transaction_id": authority}

	var payment entity.Payment
	err := d.collection.FindOne(ctx, filter).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Payment{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("payment.not_found").
				WithMeta(map[string]any{"authority": authority})
		}
		return entity.Payment{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return payment, nil
}

func (d DB) UpdatePaymentSuccess(ctx context.Context, paymentID string, refID string, cardPan string) error {
	const op = "mongodbpayment.UpdatePaymentSuccess"

	filter := bson.M{"_id": paymentID}

	update := bson.M{
		"$set": bson.M{
			"status":             entity.PaymentStatusSuccessful,
			"tracking_code":      refID,
			"masked_card_number": cardPan,
			"updated_at":         time.Now(),
		},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("payment.not_found")
	}

	return nil
}

func (d DB) UpdatePaymentFailed(ctx context.Context, paymentID string) error {
	const op = "mongodbpayment.UpdatePaymentFailed"

	filter := bson.M{"_id": paymentID}

	update := bson.M{
		"$set": bson.M{
			"status":     entity.PaymentStatusFailed,
			"updated_at": time.Now(),
		},
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("payment.not_found")
	}

	return nil
}

func (d DB) GetPaymentsByOrderID(ctx context.Context, orderID string) ([]entity.Payment, error) {
	const op = "mongodbpayment.GetPaymentsByOrderID"

	filter := bson.M{"order_id": orderID}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	defer cursor.Close(ctx)

	var payments []entity.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}

	return payments, nil
}

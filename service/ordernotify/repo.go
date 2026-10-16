package ordernotifyservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	mongodborder "github.com/abolfazlnorzad/steve/repository/mongo/order"
)

// NotifyingOrderRepo decorates the Mongo order repo so every status change —
// from the backoffice, payment verify, partner handover, QR delivery, … —
// notifies the customer without each call site having to remember to.
type NotifyingOrderRepo struct {
	*mongodborder.DB
	notifier *Service
}

func WrapOrderRepo(db *mongodborder.DB) *NotifyingOrderRepo {
	return &NotifyingOrderRepo{DB: db}
}

// SetNotifier breaks the construction cycle: the notifier itself needs the
// order repo to claim thank-you coupons.
func (r *NotifyingOrderRepo) SetNotifier(n Service) {
	r.notifier = &n
}

func (r *NotifyingOrderRepo) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	before, gErr := r.DB.GetOrderByID(ctx, orderID)

	if err := r.DB.UpdateOrderStatus(ctx, orderID, status); err != nil {
		return err
	}

	if gErr == nil && r.notifier != nil {
		r.notifier.StatusChanged(before, status)
	}
	return nil
}

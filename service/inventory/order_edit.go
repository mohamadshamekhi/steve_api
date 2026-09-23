package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// ReplaceOrderSale rewrites an edited order's sale lines so revenue, cost
// and units in the ledger match the corrected order. The old lines are
// dropped rather than reversed: they were never true, and the order itself
// keeps the before/after audit trail. The new lines keep the order's
// original date so the sale stays in the period it happened in.
//
// Only valid for an order with no return lines (cancel/sales return), which
// the order edit guarantees before calling.
func (s BackofficeService) ReplaceOrderSale(ctx context.Context, order entity.Order) error {
	const op = "inventoryservice.ReplaceOrderSale"

	if err := s.repo.DeleteByRef(ctx, entity.StockRefOrder, order.ID); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	if err := s.RecordOrderSale(ctx, order); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

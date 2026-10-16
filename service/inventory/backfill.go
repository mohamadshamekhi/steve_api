package inventoryservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// BackfillOrderSales writes the missing sale lines for orders registered
// before the sales ledger existed (or whose ledger write failed), so the
// profit reports cover them. An order that already has sale lines is left
// alone, so running it twice is safe.
//
// Cost comes from the order item's recorded cost when present, otherwise
// from the variant's current average cost (the best estimate available for
// an order older than the ledger).
func (s BackofficeService) BackfillOrderSales(ctx context.Context, from, to time.Time) (inventoryparam.BackfillSalesResponse, error) {
	const op = "inventoryservice.BackfillOrderSales"

	if s.orderProvider == nil {
		return inventoryparam.BackfillSalesResponse{}, richerror.New(op).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	orders, err := s.orderProvider.SoldOrdersBetween(ctx, from, to)
	if err != nil {
		return inventoryparam.BackfillSalesResponse{}, richerror.New(op).WithWarpError(err)
	}

	ids := make([]string, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, o.ID)
	}
	booked, err := s.repo.OrderIDsWithSales(ctx, ids)
	if err != nil {
		return inventoryparam.BackfillSalesResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := inventoryparam.BackfillSalesResponse{Scanned: len(orders)}
	for _, o := range orders {
		if _, ok := booked[o.ID]; ok || !hasItems(o) {
			resp.Skipped++
			continue
		}
		if err := s.RecordOrderSale(ctx, o); err != nil {
			logger.L().Error("backfill: failed to record order sale", "order_id", o.ID, "err", err)
			resp.Skipped++
			continue
		}
		resp.Recorded++
	}
	return resp, nil
}

func hasItems(o entity.Order) bool { return len(o.Items) > 0 }

package salesreturnservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	salesreturnparam "github.com/abolfazlnorzad/steve/param/salesreturn"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// Create books a customer return against an order: it values the returned
// goods at what they were actually sold for, puts the sellable ones back on
// the shelf, writes the ledger lines, and — when the money has already been
// handed back — takes the refund out of the cashbox.
func (s BackofficeService) Create(ctx context.Context, req salesreturnparam.CreateSalesReturnRequest, createdBy string) (salesreturnparam.CreateSalesReturnResponse, error) {
	const op = "salesreturnservice.Create"

	order, err := s.orderProvider.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).WithWarpError(err)
	}

	returned, err := s.returnedQuantities(ctx, order.ID)
	if err != nil {
		return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).WithWarpError(err)
	}

	items := make([]entity.SalesReturnItem, 0, len(req.Items))
	orderItems := make([]entity.OrderItem, 0, len(req.Items))

	for _, reqItem := range req.Items {
		orderItem, ok := findOrderItem(order, reqItem.VariantID)
		if !ok {
			return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("sales_return.item_not_in_order").
				WithMeta(map[string]any{"variant_id": reqItem.VariantID})
		}

		// More cannot come back than went out, counting returns already
		// booked against this order.
		remaining := uint64(orderItem.Quantity) - returned[reqItem.VariantID]
		if reqItem.Quantity > remaining {
			return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("sales_return.quantity_exceeds_sold").
				WithMeta(map[string]any{"variant_id": reqItem.VariantID, "remaining": remaining})
		}

		item := entity.SalesReturnItem{
			ProductID:         orderItem.ProductID,
			VariantID:         orderItem.VariantID,
			ProductName:       orderItem.ProductName,
			VariantAttributes: orderItem.VariantAttributes,
			Quantity:          reqItem.Quantity,
			UnitPrice:         netUnitPrice(orderItem),
			UnitCost:          uint64(orderItem.CostAtTime),
			Restock:           restockable(orderItem),
		}
		// An explicit false always wins: damaged goods come back into the
		// books but never back onto the shelf.
		if reqItem.Restock != nil && !*reqItem.Restock {
			item.Restock = false
		}

		items = append(items, item)
		orderItems = append(orderItems, orderItem)
	}

	r := entity.SalesReturn{
		OrderID:      order.ID,
		UserID:       order.UserID,
		Items:        items,
		RefundMethod: req.RefundMethod,
		IsRefunded:   req.IsRefunded == nil || *req.IsRefunded,
		Reason:       req.Reason,
		Note:         req.Note,
		CreatedBy:    createdBy,
	}
	if req.ReturnedAt != nil {
		r.ReturnedAt = *req.ReturnedAt
	}
	r.BeforeCreate()

	r.RefundAmount = r.ItemsAmount
	if req.RefundAmount != nil {
		// A refund larger than what the goods were sold for is refused rather
		// than clamped: it means a typo, not a decision.
		if *req.RefundAmount > r.ItemsAmount {
			return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("sales_return.refund_exceeds_items_amount").
				WithMeta(map[string]any{"items_amount": r.ItemsAmount})
		}
		r.RefundAmount = *req.RefundAmount
	}

	if err := s.repo.Create(ctx, r); err != nil {
		return salesreturnparam.CreateSalesReturnResponse{}, richerror.New(op).WithWarpError(err)
	}

	s.restock(ctx, order, r, orderItems)
	if r.IsRefunded {
		s.recordRefund(ctx, r)
	}
	return salesreturnparam.CreateSalesReturnResponse{SalesReturn: r}, nil
}

// restock puts the sellable units back and writes their ledger lines. The
// return itself is already saved at this point, so a stock write that fails
// is logged rather than swallowing the whole operation — the admin can fix
// the count with an انبارگردانی adjustment.
func (s BackofficeService) restock(ctx context.Context, order entity.Order, r entity.SalesReturn, orderItems []entity.OrderItem) {
	for i, item := range r.Items {
		if !item.Restock || item.Quantity == 0 || item.ProductID == "" {
			continue
		}

		if err := s.productProvider.ApplyStockDelta(ctx, item.ProductID, item.VariantID, int64(item.Quantity), nil); err != nil {
			logger.L().Error("failed to restock returned item", "return_id", r.ID, "variant_id", item.VariantID, "err", err)
			continue
		}

		// RecordOrderReturn reads the quantity off the order line, so it gets
		// a copy narrowed to what actually came back.
		ledgerItem := orderItems[i]
		ledgerItem.Quantity = uint(item.Quantity)
		ledgerItem.StockQuantity = uint(item.Quantity)
		ledgerItem.ExternalQuantity = 0
		ledgerItem.CostAtTime = uint(item.UnitCost)

		if err := s.inventory.RecordOrderReturn(ctx, order, ledgerItem); err != nil {
			logger.L().Error("failed to record return movement", "return_id", r.ID, "variant_id", item.VariantID, "err", err)
		}
	}
}

// recordRefund books the money leaving the till.
func (s BackofficeService) recordRefund(ctx context.Context, r entity.SalesReturn) {
	if s.cash == nil || r.RefundAmount == 0 {
		return
	}

	t := entity.CashTransaction{
		Direction:   entity.CashOut,
		Method:      r.RefundMethod,
		Amount:      r.RefundAmount,
		Description: "استرداد وجه بابت مرجوعی سفارش",
		RefType:     entity.CashRefReturn,
		RefID:       r.ID,
		OccurredAt:  r.ReturnedAt,
		CreatedBy:   r.CreatedBy,
	}
	if err := s.cash.Record(ctx, t); err != nil {
		logger.L().Error("failed to record refund in cashbox", "return_id", r.ID, "err", err)
	}
}

// MarkRefunded settles a return whose money was handed back later.
func (s BackofficeService) MarkRefunded(ctx context.Context, req salesreturnparam.MarkRefundedRequest) (salesreturnparam.MarkRefundedResponse, error) {
	const op = "salesreturnservice.MarkRefunded"

	r, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return salesreturnparam.MarkRefundedResponse{}, richerror.New(op).WithWarpError(err)
	}
	if r.IsRefunded {
		return salesreturnparam.MarkRefundedResponse{SalesReturn: r}, nil
	}

	r.IsRefunded = true
	if req.RefundMethod != "" {
		r.RefundMethod = req.RefundMethod
	}
	if err := s.repo.Update(ctx, r); err != nil {
		return salesreturnparam.MarkRefundedResponse{}, richerror.New(op).WithWarpError(err)
	}

	s.recordRefund(ctx, r)
	return salesreturnparam.MarkRefundedResponse{SalesReturn: r}, nil
}

func (s BackofficeService) Get(ctx context.Context, id string) (salesreturnparam.GetSalesReturnResponse, error) {
	const op = "salesreturnservice.Get"

	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return salesreturnparam.GetSalesReturnResponse{}, richerror.New(op).WithWarpError(err)
	}
	return salesreturnparam.GetSalesReturnResponse{SalesReturn: r}, nil
}

func (s BackofficeService) List(ctx context.Context, req salesreturnparam.ListSalesReturnsRequest) (salesreturnparam.ListSalesReturnsResponse, error) {
	const op = "salesreturnservice.List"

	returns, total, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return salesreturnparam.ListSalesReturnsResponse{}, richerror.New(op).WithWarpError(err)
	}
	summary, err := s.repo.Summary(ctx, req.Filter)
	if err != nil {
		return salesreturnparam.ListSalesReturnsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return salesreturnparam.ListSalesReturnsResponse{
		SalesReturns: returns,
		Summary:      summary,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

// Summary is the port the profit and loss statement reads.
func (s BackofficeService) Summary(ctx context.Context, f salesreturnparam.SalesReturnFilter) (salesreturnparam.SalesReturnSummary, error) {
	const op = "salesreturnservice.Summary"

	summary, err := s.repo.Summary(ctx, f)
	if err != nil {
		return salesreturnparam.SalesReturnSummary{}, richerror.New(op).WithWarpError(err)
	}
	return summary, nil
}

// Export returns the filtered rows for a spreadsheet download.
func (s BackofficeService) Export(ctx context.Context, f salesreturnparam.SalesReturnFilter) ([]entity.SalesReturn, error) {
	const op = "salesreturnservice.Export"

	returns, err := s.repo.AllForExport(ctx, f)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	return returns, nil
}

// returnedQuantities sums what has already come back per variant.
func (s BackofficeService) returnedQuantities(ctx context.Context, orderID string) (map[string]uint64, error) {
	previous, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	returned := make(map[string]uint64)
	for _, r := range previous {
		for _, item := range r.Items {
			returned[item.VariantID] += item.Quantity
		}
	}
	return returned, nil
}

func findOrderItem(order entity.Order, variantID string) (entity.OrderItem, bool) {
	for _, item := range order.Items {
		if item.VariantID == variantID {
			return item, true
		}
	}
	return entity.OrderItem{}, false
}

// netUnitPrice is what the customer actually paid per unit, after the
// campaign and manual discounts that were applied to that line — refunding
// the list price would hand back money that was never taken.
func netUnitPrice(item entity.OrderItem) uint64 {
	if item.Quantity == 0 {
		return 0
	}
	gross := uint64(item.PriceAtTime) * uint64(item.Quantity)
	discounts := uint64(item.DiscountAmount) + uint64(item.ManualDiscount)
	if discounts >= gross {
		return 0
	}
	return (gross - discounts) / uint64(item.Quantity)
}

// restockable reports whether the units can go back on the shelf: a service
// has no shelf, a hand-typed counter line has no catalog entry, and an
// externally supplied unit never entered the warehouse.
func restockable(item entity.OrderItem) bool {
	if item.IsCustom || item.ProductID == "" {
		return false
	}
	if item.GetKind() == entity.ProductKindService {
		return false
	}
	if item.Supply == entity.SupplySourceExternal {
		return false
	}
	return true
}

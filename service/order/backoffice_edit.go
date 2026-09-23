package orderservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
)

// OrderEditor persists a corrected order with optimistic locking.
type OrderEditor interface {
	ApplyEdit(ctx context.Context, order entity.Order, prevUpdatedAt time.Time, edit entity.OrderEdit) error
}

// CatalogReader resolves products/variants for new or changed lines.
type CatalogReader interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
}

// SalesReturnFinder tells whether an order already has returns booked
// against it; such an order's lines can no longer be rewritten safely.
type SalesReturnFinder interface {
	FindByOrderID(ctx context.Context, orderID string) ([]entity.SalesReturn, error)
}

// OrderLedgerRewriter replaces an edited order's sale lines in the ledger.
type OrderLedgerRewriter interface {
	ReplaceOrderSale(ctx context.Context, order entity.Order) error
}

type OrderEditDeps struct {
	Orders  OrderEditor
	Catalog CatalogReader
	Returns SalesReturnFinder
	Ledger  OrderLedgerRewriter
}

// SetOrderEditDeps wires the order-correction flow in after construction.
func (s *BackofficeService) SetOrderEditDeps(d OrderEditDeps) {
	s.edit = d
}

// stockKey identifies one catalog variant whose cached stock moves.
type stockKey struct{ productID, variantID string }

// EditOrder lets a super admin correct a mistakenly registered order.
//
// When the order has already taken stock (see entity.Order.StockApplied),
// only the difference is moved: extra units come out of stock (or are marked
// externally sourced when the shelf is short), removed units go back, and the
// order's ledger lines are rewritten. A change in the final amount can be
// booked in the cashbox. Mongo here is standalone (no multi-document
// transactions), so stock moves are compensated if a later step fails.
func (s BackofficeService) EditOrder(ctx context.Context, req orderparam.EditOrderRequest) (orderparam.EditOrderResponse, error) {
	const op = "orderservice.EditOrder"

	if s.edit.Orders == nil || s.edit.Catalog == nil {
		return orderparam.EditOrderResponse{}, richerror.New(op).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	order, err := s.repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	prev := order
	before := order.Snapshot()
	applied := order.StockApplied()
	var changes []string

	if req.UserID != nil && *req.UserID != order.UserID {
		users, uErr := s.userProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: []string{*req.UserID}})
		if uErr != nil || len(users.Users) == 0 {
			return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(uErr).
				WithKind(richerror.KindNotFound).WithMessageKey("user.not_found")
		}
		order.UserID = *req.UserID
		changes = append(changes, "user_id")
	}
	if req.ShippingAddress != nil && *req.ShippingAddress != order.ShippingAddress {
		order.ShippingAddress = *req.ShippingAddress
		changes = append(changes, "shipping_address")
	}
	if req.ShippingAmount != nil && *req.ShippingAmount != order.ShippingAmount {
		order.ShippingAmount = *req.ShippingAmount
		changes = append(changes, "shipping_amount")
	}
	if req.ManualDiscountAmount != nil && *req.ManualDiscountAmount != order.ManualDiscountAmount {
		order.ManualDiscountAmount = *req.ManualDiscountAmount
		changes = append(changes, "manual_discount_amount")
	}

	deltas := map[stockKey]int64{}
	itemsChanged := false
	if req.Items != nil {
		newItems, d, iErr := s.rebuildItems(ctx, order, req.Items, applied)
		if iErr != nil {
			return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		if !sameItems(order.Items, newItems) {
			order.Items = newItems
			deltas = d
			itemsChanged = true
			changes = append(changes, "items")
		}
	}

	if len(changes) == 0 {
		return orderparam.EditOrderResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("order.edit_no_changes")
	}

	if err := recomputeTotals(&order, prev); err != nil {
		return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	cashDiff := int64(order.FinalAmount) - int64(prev.FinalAmount)
	recordCash := order.PurchaseType == entity.PurchaseTypeInPerson &&
		order.Status == entity.OrderStatusDelivered && !order.IsPartnerOrder
	if req.RecordCashDifference != nil {
		recordCash = *req.RecordCashDifference
	}
	if !recordCash || s.cash == nil {
		cashDiff = 0
	}

	now := time.Now()
	edit := entity.OrderEdit{
		ID:             ulid.Make().String(),
		EditedBy:       req.EditedBy,
		EditedAt:       now,
		Reason:         req.Reason,
		Changes:        changes,
		Before:         before,
		After:          order.Snapshot(),
		CashDifference: cashDiff,
	}
	order.UpdatedAt = now

	undo, sErr := s.applyStockDeltas(ctx, deltas)
	if sErr != nil {
		return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(sErr)
	}
	if err := s.edit.Orders.ApplyEdit(ctx, order, prev.UpdatedAt, edit); err != nil {
		undo()
		return orderparam.EditOrderResponse{}, richerror.New(op).WithWarpError(err)
	}

	if applied && itemsChanged && s.edit.Ledger != nil {
		if err := s.edit.Ledger.ReplaceOrderSale(ctx, order); err != nil {
			logger.L().Error("order edit: failed to rewrite ledger", "order_id", order.ID, "err", err)
		}
	}
	if cashDiff != 0 {
		s.recordEditCash(ctx, order, edit, req.PaymentMethod)
	}

	order.Edits = append(order.Edits, edit)
	return orderparam.EditOrderResponse{Order: order, Edit: edit}, nil
}

// rebuildItems turns the requested lines into order items and works out how
// much cached stock each catalog variant must move (negative = out of stock).
func (s BackofficeService) rebuildItems(ctx context.Context, order entity.Order, reqItems []orderparam.EditOrderItem, applied bool) ([]entity.OrderItem, map[stockKey]int64, error) {
	const op = "orderservice.rebuildItems"

	if order.Status == entity.OrderStatusCancelled || order.Status == entity.OrderStatusRefunded {
		return nil, nil, richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessageKey("order.edit_items_not_allowed")
	}
	if len(reqItems) == 0 {
		return nil, nil, richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("order.edit_empty_items")
	}
	if s.edit.Returns != nil {
		returns, err := s.edit.Returns.FindByOrderID(ctx, order.ID)
		if err != nil {
			return nil, nil, richerror.New(op).WithWarpError(err)
		}
		if len(returns) > 0 {
			return nil, nil, richerror.New(op).WithKind(richerror.KindBadRequest).
				WithMessageKey("order.edit_has_returns")
		}
	}

	oldByVariant := make(map[string]entity.OrderItem, len(order.Items))
	for _, it := range order.Items {
		oldByVariant[it.VariantID] = it
	}

	products := map[string]entity.Product{}
	loadProduct := func(id string) (entity.Product, error) {
		if p, ok := products[id]; ok {
			return p, nil
		}
		p, err := s.edit.Catalog.FindByID(ctx, id)
		if err != nil {
			return entity.Product{}, err
		}
		products[id] = p
		return p, nil
	}

	deltas := map[stockKey]int64{}
	seen := map[string]bool{}
	items := make([]entity.OrderItem, 0, len(reqItems))

	for _, it := range reqItems {
		if it.VariantID != "" {
			if seen[it.VariantID] {
				return nil, nil, richerror.New(op).WithKind(richerror.KindInvalid).
					WithMessageKey("order.edit_duplicate_line").WithMeta(map[string]any{"variant_id": it.VariantID})
			}
			seen[it.VariantID] = true
		}
		old, exists := oldByVariant[it.VariantID]
		if it.VariantID == "" {
			exists = false
		}

		var line entity.OrderItem
		if it.IsCustom {
			line = customLine(it, old, exists)
		} else {
			if exists && old.IsCustom {
				return nil, nil, richerror.New(op).WithKind(richerror.KindInvalid).
					WithMessageKey("order.edit_line_kind_mismatch").WithMeta(map[string]any{"variant_id": it.VariantID})
			}
			productID := it.ProductID
			if exists {
				productID = old.ProductID
			}
			if productID == "" {
				return nil, nil, richerror.New(op).WithKind(richerror.KindInvalid).
					WithMessageKey("validation.required").WithMeta(map[string]any{"field": "product_id", "variant_id": it.VariantID})
			}
			product, err := loadProduct(productID)
			if err != nil {
				return nil, nil, richerror.New(op).WithWarpError(err)
			}
			variant, ok := variantOf(product, it.VariantID)
			if !ok {
				return nil, nil, richerror.New(op).WithKind(richerror.KindNotFound).
					WithMessageKey("product.variant_not_found").
					WithMeta(map[string]any{"product_id": productID, "variant_id": it.VariantID})
			}

			if exists {
				line = old
			} else {
				line = entity.OrderItem{
					ProductID:         product.ID,
					VariantID:         variant.ID,
					ProductName:       product.Name,
					VariantAttributes: variant.VariantAttribute,
					PriceAtTime:       uint(variant.Price),
				}
			}
			if it.UnitPrice != nil {
				line.PriceAtTime = *it.UnitPrice
			}

			key := stockKey{product.ID, variant.ID}
			switch {
			case !exists:
				sourceFresh(&line, variant, it.Quantity)
				if applied {
					deltas[key] -= int64(line.StockQuantity)
				}
			case it.Quantity != old.Quantity && applied:
				deltas[key] += resourceDelta(&line, old, variant, it.Quantity)
			case it.Quantity != old.Quantity:
				sourceFresh(&line, variant, it.Quantity)
			}
			line.Quantity = it.Quantity
		}

		line.DiscountAmount, line.ManualDiscount = lineDiscount(it, old, exists)
		if uint64(line.DiscountAmount) > uint64(line.PriceAtTime)*uint64(line.Quantity) {
			return nil, nil, richerror.New(op).WithKind(richerror.KindInvalid).
				WithMessageKey("order.edit_discount_exceeds_line").WithMeta(map[string]any{"variant_id": line.VariantID})
		}
		items = append(items, line)
	}

	// Lines dropped from the order give their stock back.
	if applied {
		for _, old := range order.Items {
			if seen[old.VariantID] || old.IsCustom || old.ProductID == "" || old.GetKind() == entity.ProductKindService {
				continue
			}
			if back := stockPart(old); back > 0 {
				deltas[stockKey{old.ProductID, old.VariantID}] += int64(back)
			}
		}
	}

	// Only tracked variants keep a stock count; the rest never moved.
	for key, d := range deltas {
		product, err := loadProduct(key.productID)
		if err != nil {
			logger.L().Error("order edit: product missing for stock delta", "product_id", key.productID, "err", err)
			delete(deltas, key)
			continue
		}
		if v, ok := variantOf(product, key.variantID); !ok || !v.IsStockTracked() || d == 0 {
			delete(deltas, key)
		}
	}

	return items, deltas, nil
}

// customLine builds a hand-typed line; it never touches catalog stock.
func customLine(it orderparam.EditOrderItem, old entity.OrderItem, exists bool) entity.OrderItem {
	line := entity.OrderItem{VariantID: it.VariantID, IsCustom: true}
	if exists && old.IsCustom {
		line = old
	}
	if line.VariantID == "" {
		line.VariantID = ulid.Make().String()
	}
	line.ProductID = ""
	line.ProductName = it.ProductName
	line.Kind = it.Kind
	if line.Kind == "" {
		line.Kind = entity.ProductKindPhysical
	}
	line.Quantity = it.Quantity
	if it.UnitPrice != nil {
		line.PriceAtTime = *it.UnitPrice
	}
	if it.UnitCost != nil {
		line.CostAtTime = *it.UnitCost
	}
	line.StockQuantity = 0
	if line.Kind == entity.ProductKindService {
		line.Supply = entity.SupplySourceService
		line.ExternalQuantity = 0
	} else {
		line.Supply = entity.SupplySourceExternal
		line.ExternalQuantity = it.Quantity
	}
	return line
}

// sourceFresh sources the whole quantity against the variant as it stands,
// exactly like checkout does.
func sourceFresh(line *entity.OrderItem, v entity.Variant, qty uint) {
	tmp := entity.CheckoutItemResult{}
	tmp.Quantity = int(qty)
	applySourcing(&tmp, v)
	line.Kind = tmp.Kind
	line.CostAtTime = tmp.CostAtTime
	line.Supply = tmp.Supply
	line.StockQuantity = tmp.StockQuantity
	line.ExternalQuantity = tmp.ExternalQuantity
	line.SupplierID = tmp.SupplierID
}

// resourceDelta changes the quantity of a line that already took stock and
// returns the stock movement (negative = more units out of stock). Extra
// units come from stock first, then external supply; removed units drop the
// external part first, since those never came off the shelf.
func resourceDelta(line *entity.OrderItem, old entity.OrderItem, v entity.Variant, qty uint) int64 {
	if old.GetKind() == entity.ProductKindService || v.GetStockMode() == entity.StockModeNone {
		line.Supply = entity.SupplySourceService
		return 0
	}
	stock, ext := stockPart(old), externalPart(old)

	var delta int64
	if qty > old.Quantity {
		fromStock, external := v.SplitSupply(uint64(qty - old.Quantity))
		stock += uint(fromStock)
		ext += uint(external)
		delta = -int64(fromStock)
	} else {
		remove := old.Quantity - qty
		takeExt := min(ext, remove)
		ext -= takeExt
		back := remove - takeExt
		stock -= back
		delta = int64(back)
	}

	line.StockQuantity, line.ExternalQuantity = stock, ext
	switch {
	case ext == 0:
		line.Supply = entity.SupplySourceStock
	case stock == 0:
		line.Supply = entity.SupplySourceExternal
	default:
		line.Supply = entity.SupplySourceMixed
	}
	return delta
}

func stockPart(i entity.OrderItem) uint {
	if i.Supply == "" {
		return i.Quantity
	}
	return i.StockQuantity
}

func externalPart(i entity.OrderItem) uint {
	if i.Supply == "" {
		return 0
	}
	return i.ExternalQuantity
}

// lineDiscount returns (total discount, manual part) for the line. An
// explicit discount is taken as given; the campaign part of the old line is
// kept (scaled) and the rest counted as manual. Without one, an existing
// line's discount scales with its quantity.
func lineDiscount(it orderparam.EditOrderItem, old entity.OrderItem, exists bool) (uint, uint) {
	var campaign, manual uint
	if exists && old.Quantity > 0 {
		campaign = scale(old.DiscountAmount-min(old.ManualDiscount, old.DiscountAmount), it.Quantity, old.Quantity)
		manual = scale(min(old.ManualDiscount, old.DiscountAmount), it.Quantity, old.Quantity)
	}
	if it.DiscountAmount == nil {
		return campaign + manual, manual
	}
	total := *it.DiscountAmount
	if campaign > total {
		campaign = total
	}
	return total, total - campaign
}

func scale(v, num, den uint) uint {
	if den == 0 {
		return 0
	}
	return uint(uint64(v) * uint64(num) / uint64(den))
}

// recomputeTotals rebuilds subtotal and final amount. Coupon/campaign money
// taken off the invoice as a whole (not sitting on any line) is carried over
// unchanged from the order as it was.
func recomputeTotals(order *entity.Order, prev entity.Order) error {
	const op = "orderservice.recomputeTotals"

	var prevLineDiscount int64
	for _, it := range prev.Items {
		prevLineDiscount += int64(it.DiscountAmount)
	}
	orderLevel := int64(prev.SubtotalAmount) + int64(prev.ShippingAmount) - int64(prev.FinalAmount) -
		prevLineDiscount - int64(prev.ManualDiscountAmount)
	if orderLevel < 0 {
		orderLevel = 0
	}

	var subtotal, lineDiscount int64
	for _, it := range order.Items {
		subtotal += int64(it.PriceAtTime) * int64(it.Quantity)
		lineDiscount += int64(it.DiscountAmount)
	}
	final := subtotal + int64(order.ShippingAmount) - lineDiscount - int64(order.ManualDiscountAmount) - orderLevel
	if final < 0 {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("order.edit_negative_total")
	}
	order.SubtotalAmount = uint(subtotal)
	order.FinalAmount = uint(final)
	return nil
}

// applyStockDeltas moves cached stock one variant at a time, taking stock
// out first so a shortage fails before anything is put back. The returned
// undo reverses whatever was applied.
func (s BackofficeService) applyStockDeltas(ctx context.Context, deltas map[stockKey]int64) (func(), error) {
	const op = "orderservice.applyStockDeltas"

	type done struct {
		key   stockKey
		delta int64
	}
	var applied []done
	move := func(key stockKey, delta int64) error {
		if delta < 0 {
			n := uint(-delta)
			_, err := s.productProvider.UpdateProductVariantQuantityAfterRegisterOrder(ctx,
				productparam.UpdateProductVariantQuantityAfterRegisterOrderRequest{Order: entity.Order{Items: []entity.OrderItem{{
					ProductID: key.productID, VariantID: key.variantID, Quantity: n,
					StockQuantity: n, Supply: entity.SupplySourceStock, Kind: entity.ProductKindPhysical,
				}}}})
			return err
		}
		_, err := s.productProvider.UpdateVariantQuantity(ctx, productparam.UpdateProductVariantQuantityAfterCancelOrderRequest{
			ProductID: key.productID, VariantID: key.variantID, Quantity: uint(delta),
		})
		return err
	}
	undo := func() {
		for i := len(applied) - 1; i >= 0; i-- {
			if err := move(applied[i].key, -applied[i].delta); err != nil {
				logger.L().Error("order edit: failed to undo stock move", "variant_id", applied[i].key.variantID, "err", err)
			}
		}
	}

	for _, wantOut := range []bool{true, false} {
		for key, d := range deltas {
			if (d < 0) != wantOut {
				continue
			}
			if err := move(key, d); err != nil {
				undo()
				return func() {}, richerror.New(op).WithWarpError(err)
			}
			applied = append(applied, done{key, d})
		}
	}
	return undo, nil
}

func (s BackofficeService) recordEditCash(ctx context.Context, order entity.Order, edit entity.OrderEdit, method entity.PaymentMethod) {
	if method == "" {
		method = entity.PaymentMethodCash
	}
	t := entity.CashTransaction{
		Direction:   entity.CashIn,
		Method:      method,
		Amount:      uint64(edit.CashDifference),
		Description: "اصلاح سفارش: " + edit.Reason,
		RefType:     entity.CashRefOrder,
		RefID:       order.ID,
		IsAuto:      true,
		CreatedBy:   edit.EditedBy,
		OccurredAt:  edit.EditedAt,
	}
	if edit.CashDifference < 0 {
		t.Direction = entity.CashOut
		t.Amount = uint64(-edit.CashDifference)
	}
	if err := s.cash.Record(ctx, t); err != nil {
		logger.L().Error("order edit: failed to record cash difference", "order_id", order.ID, "err", err)
	}
}

func variantOf(p entity.Product, variantID string) (entity.Variant, bool) {
	for _, v := range p.Variants {
		if v.ID == variantID {
			return v, true
		}
	}
	return entity.Variant{}, false
}

func sameItems(a, b []entity.OrderItem) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.VariantID != y.VariantID || x.Quantity != y.Quantity || x.PriceAtTime != y.PriceAtTime ||
			x.DiscountAmount != y.DiscountAmount || x.ManualDiscount != y.ManualDiscount ||
			x.ProductName != y.ProductName || x.CostAtTime != y.CostAtTime || x.Kind != y.Kind {
			return false
		}
	}
	return true
}

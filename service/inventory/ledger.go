package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// RecordOrderSale appends one ledger line per order item. Stock itself is
// already decremented by the product service at this point, so this function
// never moves quantities — it records what was sold, at what cost, and how it
// was sourced.
//
// A failure here must not fail the order: the sale is a fact, the ledger line
// is bookkeeping. Errors are logged and the remaining items still get written.
func (s BackofficeService) RecordOrderSale(ctx context.Context, order entity.Order) error {
	const op = "inventoryservice.RecordOrderSale"

	movements := make([]entity.StockMovement, 0, len(order.Items))

	for _, item := range order.Items {
		// A hand-typed counter line (repair fee, Apple ID, a part bought from
		// another shop) has no catalog product to look up, but it is still
		// revenue and still carries a cost, so it must reach the ledger.
		if item.IsCustom || item.ProductID == "" {
			t := entity.StockMovementExternalSupply
			if item.GetKind() == entity.ProductKindService {
				t = entity.StockMovementServiceSale
			}
			movements = append(movements, s.buildSaleMovement(
				t, order, item, entity.Product{}, entity.Variant{}, int64(item.Quantity)))
			continue
		}

		product, err := s.productProvider.FindByID(ctx, item.ProductID)
		if err != nil {
			logger.L().Error("inventory: product missing for sale ledger - product id : " + item.ProductID)
			continue
		}
		variant, _ := findVariant(product, item.VariantID)

		// A service sale, an external supply and a sale out of stock are all
		// revenue; they differ only in which ledger type they carry.
		stockQty := int64(item.StockQuantity)
		externalQty := int64(item.ExternalQuantity)
		if item.Supply == "" {
			// Order placed before the feature existed: treat it as fully
			// stock-served so old rows still report sensible revenue.
			stockQty = int64(item.Quantity)
		}

		if item.GetKind() == entity.ProductKindService {
			movements = append(movements, s.buildSaleMovement(
				entity.StockMovementServiceSale, order, item, product, variant, int64(item.Quantity)))
			continue
		}
		if stockQty > 0 {
			movements = append(movements, s.buildSaleMovement(
				entity.StockMovementSaleOut, order, item, product, variant, stockQty))
		}
		if externalQty > 0 {
			movements = append(movements, s.buildSaleMovement(
				entity.StockMovementExternalSupply, order, item, product, variant, externalQty))
		}
	}

	if len(movements) == 0 {
		return nil
	}
	if err := s.repo.CreateMany(ctx, movements); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

// buildSaleMovement assembles one negative-quantity ledger line.
func (s BackofficeService) buildSaleMovement(
	t entity.StockMovementType,
	order entity.Order,
	item entity.OrderItem,
	product entity.Product,
	variant entity.Variant,
	qty int64,
) entity.StockMovement {
	unitCost := uint64(item.CostAtTime)
	if unitCost == 0 {
		unitCost = variant.UnitCost()
	}

	m := entity.StockMovement{
		ProductID:         item.ProductID,
		VariantID:         item.VariantID,
		ProductName:       item.ProductName,
		Kind:              item.GetKind(),
		Type:              t,
		VariantAttributes: item.VariantAttributes,
		Quantity:          -qty,
		UnitCost:          unitCost,
		UnitPrice:         uint64(item.PriceAtTime),
		RefType:           entity.StockRefOrder,
		RefID:             order.ID,
		SupplierID:        item.SupplierID,
		QuantityAfter:     int64(variant.Quantity),
		AvgCostAfter:      variant.AvgCost,
		CreatedAt:         order.CreatedAt,
	}
	if m.ProductName == "" {
		m.ProductName = product.Name
	}
	if m.SupplierID == "" && t == entity.StockMovementExternalSupply {
		m.SupplierID = variant.SupplierID
	}
	m.BeforeCreate()
	return m
}

// RecordOrderReturn books units coming back from a cancelled or refunded
// order. Only the part that originally left the warehouse returns to it; the
// externally-supplied part never entered stock, so it is reversed as a
// cost/revenue correction only.
func (s BackofficeService) RecordOrderReturn(ctx context.Context, order entity.Order, item entity.OrderItem) error {
	const op = "inventoryservice.RecordOrderReturn"

	if item.GetKind() == entity.ProductKindService {
		return nil
	}
	// A hand-typed line never entered stock, so there is nothing to put back.
	if item.IsCustom || item.ProductID == "" {
		return nil
	}

	product, err := s.productProvider.FindByID(ctx, item.ProductID)
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	variant, ok := findVariant(product, item.VariantID)
	if !ok {
		return nil
	}

	qty := int64(item.StockQuantity)
	if item.Supply == "" {
		qty = int64(item.Quantity)
	}
	if qty == 0 {
		return nil
	}

	unitCost := uint64(item.CostAtTime)
	if unitCost == 0 {
		unitCost = variant.UnitCost()
	}

	m := entity.StockMovement{
		ProductID:         item.ProductID,
		VariantID:         item.VariantID,
		ProductName:       item.ProductName,
		Kind:              item.GetKind(),
		Type:              entity.StockMovementReturnIn,
		VariantAttributes: item.VariantAttributes,
		Quantity:          qty,
		UnitCost:          unitCost,
		RefType:           entity.StockRefOrder,
		RefID:             order.ID,
		QuantityAfter:     int64(variant.Quantity) + qty,
		AvgCostAfter:      variant.AvgCost,
		Note:              "برگشت کالا به انبار بابت لغو/مرجوعی سفارش",
	}
	m.BeforeCreate()

	if err := s.repo.Create(ctx, m); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

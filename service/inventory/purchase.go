package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// ApplyPurchase books a confirmed purchase invoice: for every item it writes
// a purchase_in ledger line and, when the goods actually enter the warehouse,
// raises the variant's quantity and folds the new price into its
// weighted-average cost.
//
// Items bought and handed straight to a customer (AddToStock=false) record
// their cost but never touch stock or the average — averaging in goods that
// were never held would distort the cost basis of what is actually on the
// shelf.
func (s BackofficeService) ApplyPurchase(ctx context.Context, p entity.Purchase) error {
	const op = "inventoryservice.ApplyPurchase"

	movements := make([]entity.StockMovement, 0, len(p.Items))

	for _, item := range p.Items {
		product, err := s.productProvider.FindByID(ctx, item.ProductID)
		if err != nil {
			return richerror.New(op).WithWarpError(err)
		}
		variant, ok := findVariant(product, item.VariantID)
		if !ok {
			return richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.variant_not_found").
				WithMeta(map[string]any{"product_id": item.ProductID, "variant_id": item.VariantID})
		}

		quantityAfter := int64(variant.Quantity)
		avgCostAfter := variant.AvgCost

		if item.AddToStock {
			avgCostAfter = entity.NewWeightedAverageCost(
				int64(variant.Quantity), variant.AvgCost, item.Quantity, item.UnitCost)
			quantityAfter = int64(variant.Quantity) + int64(item.Quantity)

			if err := s.productProvider.ApplyStockDelta(ctx, item.ProductID, item.VariantID,
				int64(item.Quantity), &avgCostAfter); err != nil {
				return richerror.New(op).WithWarpError(err)
			}
		}

		m := entity.StockMovement{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			ProductName:       item.ProductName,
			Kind:              product.GetKind(),
			Type:              entity.StockMovementPurchaseIn,
			VariantAttributes: item.VariantAttributes,
			Quantity:          int64(item.Quantity),
			UnitCost:          item.UnitCost,
			RefType:           entity.StockRefPurchase,
			RefID:             p.ID,
			SupplierID:        p.SupplierID,
			QuantityAfter:     quantityAfter,
			AvgCostAfter:      avgCostAfter,
			CreatedAt:         p.PurchasedAt,
		}
		if m.ProductName == "" {
			m.ProductName = product.Name
		}
		if !item.AddToStock {
			m.Note = "خرید بدون ورود به انبار (تحویل مستقیم به مشتری)"
		}
		m.BeforeCreate()
		movements = append(movements, m)
	}

	if err := s.repo.CreateMany(ctx, movements); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

// ReversePurchase undoes a confirmed invoice that is being cancelled: stock
// that was added is taken back out and the ledger lines are removed.
//
// The weighted average is deliberately left as-is. It cannot be un-mixed once
// later purchases and sales have folded into it, and silently guessing a
// previous value would corrupt the cost basis more than leaving it.
func (s BackofficeService) ReversePurchase(ctx context.Context, p entity.Purchase) error {
	const op = "inventoryservice.ReversePurchase"

	for _, item := range p.Items {
		if !item.AddToStock {
			continue
		}
		if err := s.productProvider.ApplyStockDelta(ctx, item.ProductID, item.VariantID,
			-int64(item.Quantity), nil); err != nil {
			return richerror.New(op).WithWarpError(err)
		}
	}

	if err := s.repo.DeleteByRef(ctx, entity.StockRefPurchase, p.ID); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

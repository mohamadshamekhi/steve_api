package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/oklog/ulid/v2"
)

// Stocktake turns a physical count into adjustment movements: for every
// counted variant the difference with the cached stock becomes one
// adjustment, all sharing one stocktake RefID.
//
// Mongo runs standalone here (no replica set), so there is no multi-document
// transaction. Instead every line is validated before anything is written,
// and if a write fails midway the stock deltas already applied are reversed,
// so a failed stocktake leaves stock as it was.
func (s BackofficeService) Stocktake(ctx context.Context, req inventoryparam.StocktakeRequest, createdBy string) (inventoryparam.StocktakeResponse, error) {
	const op = "inventoryservice.Stocktake"

	stocktakeID := ulid.Make().String()
	resp := inventoryparam.StocktakeResponse{StocktakeID: stocktakeID, Movements: []entity.StockMovement{}}

	type line struct {
		product entity.Product
		variant entity.Variant
		delta   int64
	}

	// --- 1. validate everything, write nothing ---
	products := map[string]entity.Product{}
	lines := make([]line, 0, len(req.Items))
	for i, item := range req.Items {
		product, ok := products[item.ProductID]
		if !ok {
			p, err := s.productProvider.FindByID(ctx, item.ProductID)
			if err != nil {
				return inventoryparam.StocktakeResponse{}, richerror.New(op).WithWarpError(err).
					WithMeta(map[string]any{"index": i, "product_id": item.ProductID})
			}
			product, products[item.ProductID] = p, p
		}

		variant, found := findVariant(product, item.VariantID)
		if !found {
			return inventoryparam.StocktakeResponse{}, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.variant_not_found").
				WithMeta(map[string]any{"index": i, "product_id": item.ProductID, "variant_id": item.VariantID})
		}
		if !variant.IsStockTracked() {
			return inventoryparam.StocktakeResponse{}, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("inventory.variant_not_stock_tracked").
				WithMeta(map[string]any{"index": i, "product_id": item.ProductID, "variant_id": item.VariantID})
		}

		delta := int64(*item.CountedQuantity) - int64(variant.Quantity)
		if delta == 0 {
			resp.UnchangedCount++
			continue
		}
		lines = append(lines, line{product: product, variant: variant, delta: delta})
	}

	// --- 2. apply stock deltas, reversing on failure ---
	applied := make([]line, 0, len(lines))
	rollback := func() {
		for _, l := range applied {
			if rErr := s.productProvider.ApplyStockDelta(context.Background(), l.product.ID, l.variant.ID, -l.delta, nil); rErr != nil {
				logger.L().Error("CRITICAL: stocktake rollback failed",
					"stocktake_id", stocktakeID, "product_id", l.product.ID, "variant_id", l.variant.ID,
					"delta", l.delta, "err", rErr)
			}
		}
	}

	movements := make([]entity.StockMovement, 0, len(lines))
	for _, l := range lines {
		// Found units are valued at the current average cost, so the average
		// itself does not move; missing units leave at that cost too.
		if err := s.productProvider.ApplyStockDelta(ctx, l.product.ID, l.variant.ID, l.delta, nil); err != nil {
			rollback()
			return inventoryparam.StocktakeResponse{}, richerror.New(op).WithWarpError(err).
				WithMeta(map[string]any{"product_id": l.product.ID, "variant_id": l.variant.ID})
		}
		applied = append(applied, l)

		unitCost := l.variant.UnitCost()
		m := entity.StockMovement{
			ProductID:         l.product.ID,
			VariantID:         l.variant.ID,
			ProductName:       l.product.Name,
			Kind:              l.product.GetKind(),
			Type:              entity.StockMovementAdjustment,
			VariantAttributes: l.variant.VariantAttribute,
			Quantity:          l.delta,
			UnitCost:          unitCost,
			RefType:           entity.StockRefStocktake,
			RefID:             stocktakeID,
			QuantityAfter:     int64(l.variant.Quantity) + l.delta,
			AvgCostAfter:      l.variant.AvgCost,
			Note:              req.Note,
			CreatedBy:         createdBy,
		}
		m.BeforeCreate()
		movements = append(movements, m)

		if l.delta > 0 {
			resp.SurplusUnits += uint64(l.delta)
		} else {
			resp.ShortageUnits += uint64(-l.delta)
		}
		resp.ValueDiff += l.delta * int64(unitCost)
	}

	// --- 3. write the ledger lines in one insert ---
	if err := s.repo.CreateMany(ctx, movements); err != nil {
		rollback()
		return inventoryparam.StocktakeResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp.AdjustedCount = len(movements)
	resp.Movements = movements
	return resp, nil
}

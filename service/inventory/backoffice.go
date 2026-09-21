package inventoryservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	inventoryparam "github.com/abolfazlnorzad/steve/param/inventory"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// ListMovements returns a page of the stock ledger.
func (s BackofficeService) ListMovements(ctx context.Context, req inventoryparam.ListMovementsRequest) (inventoryparam.ListMovementsResponse, error) {
	const op = "inventoryservice.ListMovements"

	movements, total, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return inventoryparam.ListMovementsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return inventoryparam.ListMovementsResponse{
		Movements: movements,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

// AdjustStock applies a manual correction (انبارگردانی) and records why.
// A positive quantity also re-averages the cost, since adding units the shop
// found is economically the same as buying them at the stated cost.
func (s BackofficeService) AdjustStock(ctx context.Context, req inventoryparam.AdjustStockRequest, createdBy string) (inventoryparam.AdjustStockResponse, error) {
	const op = "inventoryservice.AdjustStock"

	if req.Quantity == 0 {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("inventory.adjustment_quantity_required")
	}

	product, err := s.productProvider.FindByID(ctx, req.ProductID)
	if err != nil {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).WithWarpError(err)
	}
	variant, ok := findVariant(product, req.VariantID)
	if !ok {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("product.variant_not_found")
	}
	if !variant.IsStockTracked() {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("inventory.variant_not_stock_tracked")
	}
	if req.Quantity < 0 && uint64(-req.Quantity) > variant.Quantity {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("inventory.adjustment_exceeds_stock").
			WithMeta(map[string]any{"quantity": variant.Quantity})
	}

	unitCost := req.UnitCost
	if unitCost == 0 {
		unitCost = variant.UnitCost()
	}

	avgCostAfter := variant.AvgCost
	var avgCostArg *uint64
	if req.Quantity > 0 && req.UnitCost > 0 {
		avgCostAfter = entity.NewWeightedAverageCost(
			int64(variant.Quantity), variant.AvgCost, uint64(req.Quantity), req.UnitCost)
		avgCostArg = &avgCostAfter
	}

	if err := s.productProvider.ApplyStockDelta(ctx, req.ProductID, req.VariantID, req.Quantity, avgCostArg); err != nil {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).WithWarpError(err)
	}

	m := entity.StockMovement{
		ProductID:         req.ProductID,
		VariantID:         req.VariantID,
		ProductName:       product.Name,
		Kind:              product.GetKind(),
		Type:              entity.StockMovementAdjustment,
		VariantAttributes: variant.VariantAttribute,
		Quantity:          req.Quantity,
		UnitCost:          unitCost,
		RefType:           entity.StockRefManual,
		QuantityAfter:     int64(variant.Quantity) + req.Quantity,
		AvgCostAfter:      avgCostAfter,
		Note:              req.Note,
		CreatedBy:         createdBy,
	}
	m.BeforeCreate()

	if err := s.repo.Create(ctx, m); err != nil {
		return inventoryparam.AdjustStockResponse{}, richerror.New(op).WithWarpError(err)
	}
	return inventoryparam.AdjustStockResponse{Movement: m}, nil
}

// SetStockConfig writes a variant's stock mode, backorder policy, default
// supplier and cost settings.
func (s BackofficeService) SetStockConfig(ctx context.Context, req inventoryparam.SetStockConfigRequest) (inventoryparam.SetStockConfigResponse, error) {
	const op = "inventoryservice.SetStockConfig"

	if err := s.productProvider.SetVariantStockConfig(ctx, req); err != nil {
		return inventoryparam.SetStockConfigResponse{}, richerror.New(op).WithWarpError(err)
	}
	return inventoryparam.SetStockConfigResponse{}, nil
}

// ProfitReport answers "چقدر سود کردم": revenue, cost basis and gross profit
// over a window, optionally broken down by day, product or supplier.
func (s BackofficeService) ProfitReport(ctx context.Context, req inventoryparam.ProfitReportRequest) (inventoryparam.ProfitReportResponse, error) {
	const op = "inventoryservice.ProfitReport"

	summary, err := s.repo.ProfitSummary(ctx, req)
	if err != nil {
		return inventoryparam.ProfitReportResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp := inventoryparam.ProfitReportResponse{Summary: summary}
	if req.GroupBy == "" {
		return resp, nil
	}

	rows, err := s.repo.ProfitBreakdown(ctx, req)
	if err != nil {
		return inventoryparam.ProfitReportResponse{}, richerror.New(op).WithWarpError(err)
	}
	resp.Rows = rows
	return resp, nil
}

// InventoryValue reports what is sitting in the warehouse and what it is
// worth at cost.
func (s BackofficeService) InventoryValue(ctx context.Context, req inventoryparam.InventoryValueRequest) (inventoryparam.InventoryValueResponse, error) {
	const op = "inventoryservice.InventoryValue"

	resp, err := s.productProvider.InventoryValue(ctx, req)
	if err != nil {
		return inventoryparam.InventoryValueResponse{}, richerror.New(op).WithWarpError(err)
	}
	return resp, nil
}

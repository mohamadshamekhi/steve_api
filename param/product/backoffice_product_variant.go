package productparam

import "github.com/abolfazlnorzad/steve/entity"

type AddVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Price     uint64 `json:"price" bson:"price"`
	Warranty  string `json:"warranty" bson:"warranty"`
	Quantity  uint64 `json:"quantity" bson:"quantity"`
	Key       string `json:"key" bson:"key"`
	Value     string `json:"value" bson:"value"`
	Label     string `json:"label,omitempty" bson:"label,omitempty"`
	Meta      string `json:"meta,omitempty" bson:"meta,omitempty"`

	// --- inventory ---
	// StockMode defaults to "tracked". Use "on_demand" for an item the shop
	// does not hold and sources per order, "none" for a service.
	StockMode      entity.StockMode `json:"stock_mode,omitempty" bson:"stock_mode,omitempty"`
	AllowBackorder bool             `json:"allow_backorder,omitempty" bson:"allow_backorder,omitempty"`
	BackorderLimit *uint64          `json:"backorder_limit,omitempty" bson:"backorder_limit,omitempty"`
	SupplierID     string           `json:"supplier_id,omitempty" bson:"supplier_id,omitempty"`
	SupplyCost     uint64           `json:"supply_cost,omitempty" bson:"supply_cost,omitempty"`
	AvgCost        uint64           `json:"avg_cost,omitempty" bson:"avg_cost,omitempty"`
}

type AddVariantResponse struct {
}

type DeleteVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	VariantID string `json:"variant_id" bson:"variant_id"`
}

type DeleteVariantResponse struct {
}

type UpdateVariantRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	VariantID string `json:"variant_id" bson:"variant_id"`

	Price    *uint64 `json:"price" bson:"price"`
	Warranty *string `json:"warranty" bson:"warranty"`
	Quantity *uint64 `json:"quantity" bson:"quantity"`
	Key      *string `json:"key" bson:"key"`
	Value    *string `json:"value" bson:"value"`
	Label    *string `json:"label,omitempty" bson:"label,omitempty"`
	Meta     *string `json:"meta,omitempty" bson:"meta,omitempty"`

	// --- inventory ---
	StockMode      *entity.StockMode `json:"stock_mode" bson:"stock_mode"`
	AllowBackorder *bool             `json:"allow_backorder" bson:"allow_backorder"`
	BackorderLimit *uint64           `json:"backorder_limit" bson:"backorder_limit"`
	SupplierID     *string           `json:"supplier_id" bson:"supplier_id"`
	SupplyCost     *uint64           `json:"supply_cost" bson:"supply_cost"`
	AvgCost        *uint64           `json:"avg_cost" bson:"avg_cost"`
}

type UpdateVariantResponse struct {
}

// BulkUpdateVariantsRequest reuses UpdateVariantRequest per item: same fields
// as the single update-variant body, plus product_id/variant_id in the item
// itself since there's no per-product URL for a bulk request.
type BulkUpdateVariantsRequest struct {
	Items []UpdateVariantRequest `json:"items"`
}

type BulkUpdateVariantItemResult struct {
	VariantID string `json:"variant_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

type BulkUpdateVariantsResponse struct {
	Results []BulkUpdateVariantItemResult `json:"results"`
}

package productparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateProductRequest struct {
	ID                      string  `json:"id"`
	Name                    *string `json:"name"`
	ENTitle                 *string `json:"en_title"`
	BrandID                 *string `json:"brand_id"`
	CategoryID              *string `json:"category_id"`
	IsSecondHand            *bool   `json:"is_second_hand"`
	RequiresFreightShipping *bool   `json:"requires_freight_shipping"`
	Description             *string `json:"description"`

	Kind          *entity.ProductKind   `json:"kind"`
	SalesChannels []entity.SalesChannel `json:"sales_channels"`
	// Price updates a service's fee without the caller having to know about
	// the variant it is stored on. Ignored for goods.
	Price *uint64 `json:"price"`

	SeoTitle       *string             `json:"seo_title"`
	SeoDescription *string             `json:"seo_description"`
	FocusKeyword   *string             `json:"focus_keyword"`
	FAQ            []entity.ProductFAQ `json:"faq"`

	// TorobURL sets the product's Torob page; "" clears it.
	TorobURL *string `json:"torob_url"`
}

type UpdateProductResponse struct {
	Product ProductCreated `json:"product"`
}

type UpdateProductVariantQuantityAfterRegisterOrderRequest struct {
	Order entity.Order
}

type UpdateProductVariantQuantityAfterRegisterOrderResponse struct {
}

type UpdateProductVariantQuantityAfterCancelOrderRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Quantity  uint   `json:"quantity"`
}

type UpdateProductVariantQuantityAfterCancelOrderResponse struct {
}

type UpdateProductSalesCountAfterRegisterOrderRequest struct {
	Order entity.Order
}

type UpdateProductSalesCountAfterRegisterOrderResponse struct{}

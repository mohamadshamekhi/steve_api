package productparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateProductRequest struct {
	ID           string  `json:"id"`
	Name         *string `json:"name"`
	ENTitle      *string `json:"en_title"`
	BrandID      *string `json:"brand_id"`
	CategoryID   *string `json:"category_id"`
	IsSecondHand *bool   `json:"is_second_hand"`
	Description  *string `json:"description"`
}

type UpdateProductResponse struct {
	Product entity.Product `json:"product"`
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

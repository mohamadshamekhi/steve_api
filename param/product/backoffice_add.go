package productparam

import "github.com/abolfazlnorzad/steve/entity"

type CreateProductRequest struct {
	Name         string `json:"name"`
	ENTitle      string `json:"en_title"`
	BrandID      string `json:"brand_id"`
	CategoryID   string `json:"category_id"`
	IsSecondHand bool   `json:"is_second_hand"`
	Description  string `json:"description"`
}
type CreateProductResponse struct {
	Product entity.Product `json:"product"`
}

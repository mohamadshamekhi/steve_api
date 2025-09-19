package productparam

import "github.com/abolfazlnorzad/steve/entity"

type GetProductCategoriesForCartItemsRequest struct {
	Items []entity.CartItem `json:"items"`
}

type GetProductCategoriesForCartItemsResponse struct {
	CategoryIDs map[string][]string `json:"category_ids"`
}

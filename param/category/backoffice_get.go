package categoryparam

import "github.com/abolfazlnorzad/steve/entity"

type GetCategoryRequest struct {
	ID string `json:"id"`
}

type GetCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}

type GetCategoriesByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetCategoriesByIDsResponse struct {
	Categories map[string]entity.Category `json:"categories"`
}

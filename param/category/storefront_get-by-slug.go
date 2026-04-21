package categoryparam

import "github.com/abolfazlnorzad/steve/entity"

type GetCategoryBySlugRequest struct {
	Slug string `json:"slug"`
}

type GetCategoryBySlugResponse struct {
	Category entity.Category `json:"category"`
}

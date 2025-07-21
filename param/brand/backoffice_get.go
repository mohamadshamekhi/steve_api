package brandparam

import "github.com/abolfazlnorzad/steve/entity"

type GetBrandRequest struct {
	ID string `json:"id"`
}

type GetBrandResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type GetBrandsByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetBrandsByIDsResponse struct {
	Brands map[string]entity.Brand `json:"brands"`
}

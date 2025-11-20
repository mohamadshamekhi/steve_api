package productparam

type GetProductIDsByCategoriesRequest struct {
	CatIDs []string
}

type GetProductIDsByCategoriesResponse struct {
	ProdsByCat map[string][]string
}

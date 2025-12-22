package productparam

type GetPriceOfProductsRequest struct {
	ProductIDs []string `json:"product_i_ds"`
	VariantIDs []string `json:"variant_id"`
}

type GetPriceOfProductsResponse struct {
	Price map[string]map[string]float64 `json:"price"`
}

package productparam

type GetProductsWithSpecialDiscountRequest struct {
}

type GetProductsWithSpecialDiscountResponse struct {
	Products []StorefrontProductList `json:"products"`
}

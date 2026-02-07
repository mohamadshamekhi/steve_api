package productparam

type UpdateProductSalesCountRequest struct {
	ProductID string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}

type UpdateProductSalesCountResponse struct{}

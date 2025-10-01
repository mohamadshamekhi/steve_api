package productparam

type StorefrontGetProductRequest struct {
	ID string `json:"id"`
}
type StorefrontGetProductResponse struct {
	Product ProductInfo `json:"product"`
}

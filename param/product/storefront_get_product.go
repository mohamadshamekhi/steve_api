package productparam

type StorefrontGetProductRequest struct {
	ID string `json:"id"`
}
type StorefrontGetProductResponse struct {
	Product ProductInfo `json:"product"`
}

type StorefrontGetProductBySlugRequest struct {
	Slug string `json:"slug"`
}
type StorefrontGetProductBySlugResponse struct {
	Product ProductInfo `json:"product"`
}

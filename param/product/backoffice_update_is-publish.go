package productparam

type PublishProductRequest struct {
	ProductID string `json:"product_id"`
}

type PublishProductResponse struct {
}

type UnPublishProductRequest struct {
	ProductID string `json:"product_id"`
}

type UnPublishProductResponse struct {
}

package productparam

type AddProductGalleryRequest struct {
	MediaID   string `json:"media_id"`
	MediaURL  string `json:"media_url"`
	IsPrimary bool   `json:"is_primary"`
	ProductID string `json:"product_id"`
	Alt       string `json:"alt"`
}

type AddProductGalleryResponse struct {
}

type SetProductGalleryAltRequest struct {
	MediaID   string `json:"media_id"`
	ProductID string `json:"product_id"`
	Alt       string `json:"alt"`
}

type SetProductGalleryAltResponse struct {
}

type DeleteProductGalleryRequest struct {
	MediaID   string `json:"media_id"`
	ProductID string `json:"product_id"`
}

type DeleteProductGalleryResponse struct {
}

type SetPrimaryProductGalleryRequest struct {
	MediaID   string `json:"media_id"`
	ProductID string `json:"product_id"`
}

type SetPrimaryProductGalleryResponse struct {
}

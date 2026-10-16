package productparam

type ResolveVariantRequest struct {
	VariantID string `param:"variant_id"`
}

type ResolveVariantResponse struct {
	ProductID    string `json:"product_id"`
	VariantID    string `json:"variant_id"`
	Slug         string `json:"slug"`
	ProductName  string `json:"product_name"`
	VariantLabel string `json:"variant_label"`
	Kind         string `json:"kind"`
	// OnWebsite tells the QR landing page whether the product page exists on
	// the storefront; in-shop-only items and unpublished products resolve too.
	OnWebsite bool `json:"on_website"`
}

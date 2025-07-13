package brandparam

type GetBrandRequest struct {
	ID string `json:"id"`
}

type GetBrandResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

package mediaparam

type GetUrlRequest struct {
	ID string `json:"id"`
}

type GetUrlResponse struct {
	URL string `json:"url"`
}

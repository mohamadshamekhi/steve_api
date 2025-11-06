package ratingparam

type UpsertRatingRequest struct {
	ProductID string  `json:"product_id"`
	UserID    string  `json:"user_id"`
	Rating    float64 `json:"rating"`
}

type UpsertRatingResponse struct {
	AverageRating float64 `json:"average_rating"`
}

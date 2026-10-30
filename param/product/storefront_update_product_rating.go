package productparam

type UpdateProductRatingRequest struct {
	ProductID     string  `json:"product_id"`
	AverageRating float64 `json:"average_rating"`
	RatingCount   int     `json:"rating_count"`
}
type UpdateProductRatingResponse struct {
}

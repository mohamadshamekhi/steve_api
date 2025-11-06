package productparam

type UpdateProductRatingRequest struct {
	ProductID     string  `json:"product_id"`
	AverageRating float64 `json:"average_rating"`
}
type UpdateProductRatingResponse struct {
}

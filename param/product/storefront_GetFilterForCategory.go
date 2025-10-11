package productparam

type GetFilterForCategoryRequest struct {
	CategoryID string `json:"category_id"`
}

type GetFilterForCategoryResponse struct {
	AggregatedFilterOption []AggregatedFilterOption
}

type AggregatedFilterOption struct {
	Key    string   `bson:"_id"`
	Values []string `bson:"values"`
}

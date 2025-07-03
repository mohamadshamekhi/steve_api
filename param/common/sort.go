package commonparam

type SortRequest struct {
	Field     string `json:"field" ,query:"field"`
	Direction string `json:"direction" ,query:"direction"`
}

package userparam

type StoreFrontAddAddressRequest struct {
	UserID   string `json:"user_id"`
	Title    string `json:"title"`
	Province string `json:"province"`
	City     string `json:"city"`
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
}

type StoreFrontAddAddressResponse struct {
}

type StoreFrontUpdateAddressRequest struct {
	UserID   string  `json:"user_id"`
	ID       string  `json:"id"`
	Title    *string `json:"title"`
	Province *string `json:"province"`
	City     *string `json:"city"`
	Street   *string `json:"street"`
	Postcode *string `json:"postcode"`
}

type StoreFrontUpdateAddressResponse struct {
}

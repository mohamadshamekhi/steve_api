package userparam

type StoreFrontUpdateProfileRequest struct {
	UserID      string  `json:"user_id"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
}

type StoreFrontUpdateProfileResponse struct {
}

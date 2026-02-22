package userparam

type UpdateUserBackofficeRequest struct {
	UserID string  `json:"user_id"`
	Name   *string `json:"name"`
}

type UpdateUserBackofficeResponse struct {
}

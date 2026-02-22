package userparam

type UpsertUserRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type UpsertUserResponse struct {
	User UserInfo `json:"user"`
}

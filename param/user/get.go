package userparam

type GetUsersRequest struct {
	ID string `json:"id"`
}

type GetUsersResponse struct {
	User UserInfo `json:"user"`
}

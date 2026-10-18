package backofficeparam

type IssueServiceTokenRequest struct {
	ID string `json:"id"`
}

type IssueServiceTokenResponse struct {
	Token string `json:"token"`
}

package backofficeauthparam

type SendLoginCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
}
type SendLoginCodeResponse struct {
}

type VerifyLoginCodeRequest struct {
	Code        string `json:"code"`
	PhoneNumber string `json:"phone_number"`
}

type VerifyLoginCodeResponse struct {
	Token string `json:"token"`
}

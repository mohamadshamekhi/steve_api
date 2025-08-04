package userauthparam

type UserSendLoginCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
}
type UserSendLoginCodeResponse struct {
}

type UserVerifyLoginCodeRequest struct {
	Code        string `json:"code"`
	PhoneNumber string `json:"phone_number"`
}

type UserVerifyLoginCodeResponse struct {
	Token string `json:"token"`
}

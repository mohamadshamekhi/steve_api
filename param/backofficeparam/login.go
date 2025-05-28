package backofficeparam

type SendLoginCodeRequest struct {
	PhoneNumber string
}
type SendLoginCodeResponse struct {
	Message string
}

type VerifyLoginCodeRequest struct {
	Code        string
	PhoneNumber string
}

type VerifyLoginCodeResponse struct {
	Token   string
	Message string
}

package backofficeparam

type GetBackofficeRequest struct {
	ID string `json:"id"`
}

type GetBackofficeResponse struct {
	Backoffice BackofficeInfo `json:"backoffice"`
}

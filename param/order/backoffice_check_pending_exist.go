package orderparam

type CheckPendingExistRequest struct {
	UserID string `json:"user_id"`
}

type CheckPendingExistResponse struct {
	Order *BackofficeOrderInfo `json:"order"`
}

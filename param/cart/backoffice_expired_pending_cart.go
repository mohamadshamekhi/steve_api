package cartparam

import "time"

type CancelExpiredPendingCartRequest struct {
	CanselExpiredPendingCart time.Duration `json:"cansel_expired_pending_cart"`
}

type CancelExpiredPendingCartResponse struct {
}

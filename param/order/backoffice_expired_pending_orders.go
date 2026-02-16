package orderparam

import "time"

type CancelExpiredPendingOrdersRequest struct {
	CanselExpiredPendingOrders time.Duration `json:"cansel_expired_pending_orders"`
}

type CancelExpiredPendingOrdersResponse struct {
}

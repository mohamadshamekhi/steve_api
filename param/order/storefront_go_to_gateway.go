package orderparam

type RegisterOrderRequest struct {
	DiscountCode string `json:"discount_code"`
	AddressID    string `json:"address_id"`
	ShippingID   string `json:"shipping_id"`
}

type RegisterOrderResponse struct {
	OrderID string `json:"order_id"`
}

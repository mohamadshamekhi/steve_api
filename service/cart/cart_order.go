package cartservice

import orderservice "github.com/abolfazlnorzad/steve/service/order"

type CartOrder struct {
	OrderSvc orderservice.StoreFrontService
}

func NewCartOrder(orderSvc orderservice.StoreFrontService) CartOrder {
	return CartOrder{
		OrderSvc: orderSvc,
	}
}

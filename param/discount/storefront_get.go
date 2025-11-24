package discountparam

import "github.com/abolfazlnorzad/steve/entity"

type GetValidDiscountsForSpecificCartRequest struct {
	DiscountCode string
	Cart         entity.Cart
}

type GetValidDiscountsForSpecificCartResponse struct {
	Discounts []entity.Discount
}

type GetSpecialDiscountRequest struct {
}
type GetSpecialDiscountResponse struct {
	Discount entity.Discount `json:"discount"`
}

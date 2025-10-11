package discountparam

import "github.com/abolfazlnorzad/steve/entity"

type GetProductIDsThatHaveActiveAutoDiscountRequest struct {
}

type GetProductIDsThatHaveActiveAutoDiscountResponse struct {
	DiscountMap map[string]entity.Discount
}

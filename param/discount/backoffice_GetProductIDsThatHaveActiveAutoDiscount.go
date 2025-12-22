package discountparam

import (
	productparam "github.com/abolfazlnorzad/steve/param/product"
)

type GetProductIDsThatHaveActiveAutoDiscountRequest struct {
}

type GetProductIDsThatHaveActiveAutoDiscountResponse struct {
	DiscountMap map[string][]productparam.ProductDiscount
}

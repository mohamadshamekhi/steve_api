package discountparam

import "github.com/abolfazlnorzad/steve/entity"

type GetSingleDiscountRequest struct {
	ID string `json:"id"`
}

type GetSingleDiscountResponse struct {
	Discount entity.Discount `json:"discount"`
}

type GetSingleDiscountByCodeRequest struct {
	Code string `json:"code"`
}

type GetSingleDiscountByCodeResponse struct {
	Discount entity.Discount `json:"discount"`
}

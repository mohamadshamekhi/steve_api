package discountservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
)

type ListResult struct {
	Discounts  []discountparam.DiscountsList
	TotalCount int64
}

type DiscountReader interface {
	FindByID(ctx context.Context, id string) (entity.Discount, error)
	FindByCode(ctx context.Context, code string) (entity.Discount, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	FindApplicableDiscounts(ctx context.Context, code string) ([]entity.Discount, error)
}

type DiscountWriter interface {
	Create(ctx context.Context, discount entity.Discount) error
	Update(ctx context.Context, discount entity.Discount) error
	UpdateUsageCount(ctx context.Context, discountID string) error
}

type DiscountReadWriter interface {
	DiscountWriter
	DiscountReader
}

package settingservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
)

type ShippingWriter interface {
	Create(ctx context.Context, shipping entity.Shipping) error
	Update(ctx context.Context, shipping entity.Shipping) error
}
type ListResult struct {
	Shipping   []entity.Shipping
	TotalCount int64
}
type ShippingReader interface {
	FindByID(ctx context.Context, id string) (entity.Shipping, error)
	GetAll(ctx context.Context, f settingparam.FilterShipping, p commonparam.PaginationRequest) (ListResult, error)
}

type ShippingReaderWriter interface {
	ShippingWriter
	ShippingReader
}


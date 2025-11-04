package contactusservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
)

type ListResult struct {
	Orders     []entity.Order
	TotalCount int64
}

type ContactUsReader interface {
	//GetAll(ctx context.Context, f orderparam.BackofficeFilterOrder, p commonparam.PaginationRequest) (ListResult, error)
}

type ContactUsWriter interface {
	Create(ctx context.Context, c entity.ContactUs) error
}

type ContactUsReadWriter interface {
	ContactUsReader
	ContactUsWriter
}

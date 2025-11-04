package contactusservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type ListResult struct {
	Contacts   []entity.ContactUs
	TotalCount int64
}

type ContactUsReader interface {
	All(ctx context.Context, p commonparam.PaginationRequest) (ListResult, error)
}

type ContactUsWriter interface {
	Create(ctx context.Context, c entity.ContactUs) error
}

type ContactUsReadWriter interface {
	ContactUsReader
	ContactUsWriter
}

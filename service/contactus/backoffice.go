package contactusservice

import (
	"context"

	contactusparam "github.com/abolfazlnorzad/steve/param/contactus"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo ContactUsReader
}

func NewBackofficeService(r ContactUsReader) BackofficeService {
	return BackofficeService{repo: r}
}

func (s BackofficeService) GetAllContactUsRequest(ctx context.Context, req contactusparam.GetAllContactRequest) (contactusparam.GetAllContactResponse, error) {
	const op = "contactusservice.GetAllContactRequest"
	resp, err := s.repo.All(ctx, req.Pagination)
	if err != nil {
		return contactusparam.GetAllContactResponse{}, richerror.New(op).WithWarpError(err)
	}

	return contactusparam.GetAllContactResponse{Contacts: resp.Contacts, Pagination: struct {
		PageSize   uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
		Total      uint `json:"total"`
	}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)}}, nil
}

package contactusservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/param/contactus"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	repo ContactUsWriter
}

func NewStorefrontService(r ContactUsWriter) StorefrontService {
	return StorefrontService{repo: r}
}

func (s StorefrontService) Create(ctx context.Context, req contactusparam.AddNewContactRequest) (contactusparam.AddNewContactResponse, error) {
	const op = "contactusservice.Create"
	c := entity.ContactUs{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Subject:     req.Subject,
		Description: req.Description,
	}
	c.BeforeCreate()
	if err := s.repo.Create(ctx, c); err != nil {
		return contactusparam.AddNewContactResponse{}, richerror.New(op).WithWarpError(err)
	}
	return contactusparam.AddNewContactResponse{}, nil
}

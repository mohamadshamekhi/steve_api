package contactusservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/param/contactus"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	adminnotifyservice "github.com/abolfazlnorzad/steve/service/adminnotify"
)

type StorefrontService struct {
	repo     ContactUsWriter
	notifier AdminNotifier
}

// AdminNotifier tells the shop's admins about a customer submission.
type AdminNotifier interface {
	Notify(e adminnotifyservice.Event)
}

// SetAdminNotifier wires the admin SMS in after construction.
func (s *StorefrontService) SetAdminNotifier(n AdminNotifier) {
	s.notifier = n
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
	if s.notifier != nil {
		detail := c.Subject
		if detail == "" {
			detail = c.Description
		}
		s.notifier.Notify(adminnotifyservice.Event{
			Title:  "پیام تماس با ما",
			Phone:  c.PhoneNumber,
			Name:   c.Email,
			Detail: detail,
		})
	}
	return contactusparam.AddNewContactResponse{}, nil
}

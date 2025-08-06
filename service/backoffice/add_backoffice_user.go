package backofficeservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) AddNewBackofficeUser(ctx context.Context, req backofficeparam.AddNewBackofficeUserRequest) (backofficeparam.AddNewBackofficeUserResponse, error) {
	const op = "backofficeservice.AddNewBackofficeUser"

	u := entity.BackOffice{
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		NationalCode: req.NationCode,
		Address:      req.Address,
		ZipCode:      req.ZipCode,
		IsActive:     *req.IsActive,
		Roles:        req.Roles,
		Permissions:  req.Permissions,
	}

	u.BeforeCreate()

	if err := s.repo.AddNewBackofficeUser(ctx, u); err != nil {
		return backofficeparam.AddNewBackofficeUserResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.AddNewBackofficeUserResponse{}, nil
}

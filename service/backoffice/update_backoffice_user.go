package backofficeservice

import (
	"context"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) UpdateBackofficeUser(ctx context.Context, req backofficeparam.UpdateBackofficeUserRequest) (backofficeparam.UpdateBackofficeUserResponse, error) {
	const op = "backofficeservice.UpdateBackofficeUser"

	u, gErr := s.repo.GetBackOfficeUserByID(ctx, req.ID)
	if gErr != nil {
		return backofficeparam.UpdateBackofficeUserResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	// todo - move this if block code to repo layer!
	if req.NationCode != nil {
		u.NationalCode = *req.NationCode
	}
	if req.PhoneNumber != nil {
		u.PhoneNumber = *req.PhoneNumber
	}
	if req.Permissions != nil {
		u.Permissions = *req.Permissions
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if req.Address != nil {
		u.Address = *req.Address
	}
	if req.Roles != nil {
		u.Roles = *req.Roles
	}
	if req.ZipCode != nil {
		u.ZipCode = *req.ZipCode
	}

	if err := s.repo.UpdateBackofficeUser(ctx, u); err != nil {
		return backofficeparam.UpdateBackofficeUserResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.UpdateBackofficeUserResponse{}, nil
}

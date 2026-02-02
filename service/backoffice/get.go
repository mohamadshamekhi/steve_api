package backofficeservice

import (
	"context"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) GetBackoffice(ctx context.Context, req backofficeparam.GetBackofficeRequest) (backofficeparam.GetBackofficeResponse, error) {
	const op = "backoffice.GetBackoffice"
	user, err := s.repo.GetBackOfficeUserByID(ctx, req.ID)
	if err != nil {
		return backofficeparam.GetBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.GetBackofficeResponse{
		Backoffice: backofficeparam.BackofficeInfo{
			ID:           user.ID,
			Name:         user.Name,
			PhoneNumber:  user.PhoneNumber,
			NationalCode: user.NationalCode,
			Address:      user.Address,
			ZipCode:      user.ZipCode,
			IsActive:     user.IsActive,
			IsSuperAdmin: user.IsSuperAdmin(),
			Roles:        user.Roles,
			Permissions:  user.Permissions,
		},
	}, nil
}

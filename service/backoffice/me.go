package backofficeservice

import (
	"context"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) GetMyProfile(ctx context.Context, req backofficeparam.MeRequest) (backofficeparam.MeResponse, error) {
	const op = "backoffice.GetMyProfile"
	user, err := s.repo.GetBackOfficeUserByID(ctx, req.ID)
	if err != nil {
		return backofficeparam.MeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.MeResponse{
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

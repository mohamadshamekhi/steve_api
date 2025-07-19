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

	if err := s.repo.UpdateBackofficeUser(ctx, u); err != nil {
		return backofficeparam.UpdateBackofficeUserResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.UpdateBackofficeUserResponse{}, nil
}

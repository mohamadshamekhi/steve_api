package backofficeservice

import (
	"context"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) Logout(ctx context.Context, req backofficeparam.LogoutRequest) (response backofficeparam.LogoutResponse, err error) {
	const op = "backofficeservice.Logout"

	if err := s.cacheProvider.AddToBlocklist(ctx, req.Kind, req.JWTID, req.ExpireAt); err != nil {
		return backofficeparam.LogoutResponse{}, richerror.New(op).WithWarpError(err)
	}

	return backofficeparam.LogoutResponse{}, nil

}

package backofficeservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) CheckBackOfficeUserIsActive(context context.Context, id string) (bool, error) {
	const op = "backofficeservice.CheckBackOfficeUserIsActive"
	bu, gErr := s.repo.GetBackOfficeUserByID(context, id)
	if gErr != nil {
		return false, richerror.New(op).WithWarpError(gErr)
	}
	
	return bu.IsActive, nil
}

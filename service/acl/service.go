package aclservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/abolfazlnorzad/steve/pkg/slice"
)

type RoleRepository interface {
	GetPermissionIDsOfRole(ctx context.Context, roleID string) ([]string, error)
}
type Service struct {
	RoleRepo RoleRepository
}

func New(roleRepo RoleRepository) Service {
	return Service{RoleRepo: roleRepo}
}

func (s Service) CheckAccess(ctx context.Context, backofficeUserID string,
	backofficeUserRoleIDs []string, backofficeUserPermissionIDs []string, permissionIDs ...string) (bool, error) {
	const op = "aclservice.CheckAccess"
	// get permissions of backoffice roles
	pIDs, gErr := s.RoleRepo.GetPermissionIDsOfRole(ctx, backofficeUserID)
	if gErr != nil {
		return false, richerror.New(op).WithWarpError(gErr)
	}
	finalPermissionIDs := []string{}

	for _, pid := range backofficeUserPermissionIDs {
		if !slice.DoesExist(finalPermissionIDs, pid) {
			finalPermissionIDs = append(finalPermissionIDs, pid)
		}
	}
	for _, pid := range pIDs {
		if !slice.DoesExist(finalPermissionIDs, pid) {
			finalPermissionIDs = append(finalPermissionIDs, pid)
		}
	}

	for _, pid := range finalPermissionIDs {
		for _, p := range permissionIDs {
			if p == pid {
				return true, nil
			}
		}
	}

	return false, nil
}

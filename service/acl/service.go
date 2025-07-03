package aclservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type RoleRepository interface {
	GetPermissionKeysOfRole(ctx context.Context, roleKeys []string) ([]string, error)
}
type Service struct {
	RoleRepo RoleRepository
}

func New(roleRepo RoleRepository) Service {
	return Service{RoleRepo: roleRepo}
}

func (s Service) CheckAccess(ctx context.Context, backofficeUserID string,
	backofficeUserRoleKeys []string, backofficeUserPermissionKeys []string, requiredPermissions ...string) (bool, error) {
	const op = "aclservice.CheckAccess"
	rolePermissionKeys, gErr := s.RoleRepo.GetPermissionKeysOfRole(ctx, backofficeUserRoleKeys)
	if gErr != nil {
		return false, richerror.New(op).WithWarpError(gErr)
	}

	userPermissionsSet := make(map[string]struct{})

	for _, k := range backofficeUserPermissionKeys {
		userPermissionsSet[k] = struct{}{}
	}

	for _, pkey := range rolePermissionKeys {
		userPermissionsSet[pkey] = struct{}{}
	}

	for _, requiredPerm := range requiredPermissions {
		if _, ok := userPermissionsSet[requiredPerm]; ok {
			return true, nil
		}
	}

	return false, nil
}

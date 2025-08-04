package aclservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type RoleRepository interface {
	GetPermissionKeysOfRole(ctx context.Context, roleKeys []string) ([]string, error)
}
type Service struct {
	RoleRepo       RoleRepository
	BackofficeRepo BackofficeRepository
}
type BackofficeRepository interface {
	GetBackOfficeUserByID(ctx context.Context, id string) (entity.BackOffice, error)
}

func New(roleRepo RoleRepository, BackofficeRepo BackofficeRepository) Service {
	return Service{RoleRepo: roleRepo, BackofficeRepo: BackofficeRepo}
}

func (s Service) CheckAccess(ctx context.Context, backofficeUserID string, requiredPermissions ...string) (bool, error) {
	const op = "aclservice.CheckAccess"
	user, err := s.BackofficeRepo.GetBackOfficeUserByID(ctx, backofficeUserID)
	if err != nil {
		return false, err
	}
	backofficeUserRoleKeys := user.Roles
	backofficeUserPermissionKeys := user.Permissions
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

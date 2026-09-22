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

// PermissionSet returns every permission the backoffice user holds: the ones
// granted directly plus the ones of each of their roles.
func (s Service) PermissionSet(ctx context.Context, backofficeUserID string) (map[string]struct{}, error) {
	const op = "aclservice.PermissionSet"
	user, err := s.BackofficeRepo.GetBackOfficeUserByID(ctx, backofficeUserID)
	if err != nil {
		return nil, err
	}
	rolePermissionKeys, gErr := s.RoleRepo.GetPermissionKeysOfRole(ctx, user.Roles)
	if gErr != nil {
		return nil, richerror.New(op).WithWarpError(gErr)
	}

	set := make(map[string]struct{}, len(user.Permissions)+len(rolePermissionKeys))
	for _, k := range user.Permissions {
		set[k] = struct{}{}
	}
	for _, k := range rolePermissionKeys {
		set[k] = struct{}{}
	}
	return set, nil
}

func (s Service) CheckAccess(ctx context.Context, backofficeUserID string, requiredPermissions ...string) (bool, error) {
	userPermissionsSet, err := s.PermissionSet(ctx, backofficeUserID)
	if err != nil {
		return false, err
	}

	for _, requiredPerm := range requiredPermissions {
		if _, ok := userPermissionsSet[requiredPerm]; ok {
			return true, nil
		}
	}

	return false, nil
}

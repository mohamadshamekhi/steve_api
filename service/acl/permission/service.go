package permissionservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	aclparam "github.com/abolfazlnorzad/steve/param/acl"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Repository interface {
	GetAllPermissions(ctx context.Context, search string) ([]entity.Permission, error)
}
type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) GetAllPermissions(ctx context.Context, req aclparam.GetAllPermissionsRequest) (aclparam.GetAllPermissionsResponse, error) {
	const op = "permissionservice.GetAllPermissions"
	permissions, err := s.repo.GetAllPermissions(ctx, req.Search)
	if err != nil {
		return aclparam.GetAllPermissionsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return aclparam.GetAllPermissionsResponse{Permissions: permissions}, nil
}

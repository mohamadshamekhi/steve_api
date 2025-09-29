package roleservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	aclparam "github.com/abolfazlnorzad/steve/param/acl"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Repository interface {
	GetAllRoles(ctx context.Context, search string) ([]entity.Role, error)
}
type Service struct {
	repo Repository
	permissionProvider
}

type permissionProvider interface {
	GetAllPermissions(ctx context.Context, req aclparam.GetAllPermissionsRequest) (aclparam.GetAllPermissionsResponse, error)
}

func New(repo Repository, provider permissionProvider) Service {
	return Service{repo: repo, permissionProvider: provider}
}

func (s Service) GetAllRoles(ctx context.Context, req aclparam.GetAllRolesRequest) (aclparam.GetAllRolesResponse, error) {
	const op = "aclhandler.GetAllRoles"

	permissions, err := s.permissionProvider.GetAllPermissions(ctx, aclparam.GetAllPermissionsRequest{
		Search: "",
	})
	if err != nil {
		return aclparam.GetAllRolesResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected)
	}

	roles, err := s.repo.GetAllRoles(ctx, req.Search)
	if err != nil {
		return aclparam.GetAllRolesResponse{}, richerror.New(op).WithWarpError(err)
	}

	roleInfos := make([]aclparam.RoleInfo, 0, len(roles))
	for _, r := range roles {
		rp := []entity.Permission{}
		for _, p := range permissions.Permissions {
			for _, key := range r.PermissionKeys {
				if key == p.Name {
					rp = append(rp, p)
				}
			}
		}
		roleInfos = append(roleInfos, aclparam.RoleInfo{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Permission:  rp,
		})
	}

	return aclparam.GetAllRolesResponse{Roles: roleInfos}, nil
}

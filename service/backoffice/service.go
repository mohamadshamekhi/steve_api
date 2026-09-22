package backofficeservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
)

type Config struct {
	Prefix string `koanf:"prefix"`
}

type ListResult struct {
	Backoffice []entity.BackOffice
	TotalCount int64
}

type Repository interface {
	GetBackOfficeUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.BackOffice, error)
	GetBackOfficeUserByID(ctx context.Context, id string) (entity.BackOffice, error)
	AddNewBackofficeUser(ctx context.Context, u entity.BackOffice) error
	UpdateBackofficeUser(context.Context, entity.BackOffice) error
	AllBackoffice(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	AllSuperUserBackOffices(ctx context.Context) (ListResult, error)
}
// RolePermissionProvider resolves the current permission keys of the given
// roles from the database, so the login token reflects the role's latest
// permissions instead of only the user's directly-assigned ones.
type RolePermissionProvider interface {
	GetPermissionKeysOfRole(ctx context.Context, roleKeys []string) ([]string, error)
}

type Service struct {
	config        Config
	authSvc       authservice.Service
	repo          Repository
	cacheProvider CacheProvider
	roleRepo      RolePermissionProvider
}

type CacheProvider interface {
	AddToBlocklist(ctx context.Context, kind string, JWTID string, ExpireAt time.Time) error
}

func New(cfg Config, authSvc authservice.Service, repo Repository, cacheProvider CacheProvider, roleRepo RolePermissionProvider) Service {
	return Service{
		config:        cfg,
		authSvc:       authSvc,
		repo:          repo,
		cacheProvider: cacheProvider,
		roleRepo:      roleRepo,
	}
}

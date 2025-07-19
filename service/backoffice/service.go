package backofficeservice

import (
	"context"
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
}
type Service struct {
	config  Config
	authSvc authservice.Service
	repo    Repository
}

func New(cfg Config, authSvc authservice.Service, repo Repository) Service {
	return Service{
		config:  cfg,
		authSvc: authSvc,
		repo:    repo,
	}
}

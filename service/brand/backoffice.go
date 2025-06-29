package brandservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo BrandReadWriter
}

func NewBackofficeService(repo BrandReadWriter) BackofficeService {
	return BackofficeService{repo: repo}
}

func (s BackofficeService) Create(ctx context.Context, req brandparam.CreateBrandRequest) (brandparam.CreateBrandResponse, error) {
	const op = "brandservice.Create"

	brand := entity.Brand{
		Name:        req.Name,
		Description: req.Description,
		MediaID:     req.MediaID,
	}
	brand.BeforeCreate()

	if err := s.repo.Create(ctx, brand); err != nil {
		return brandparam.CreateBrandResponse{}, richerror.New(op).WithWarpError(err)
	}

	return brandparam.CreateBrandResponse{}, nil
}

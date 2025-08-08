package variantservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	variantparam "github.com/abolfazlnorzad/steve/param/variant"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, variant entity.Variant) error
	Update(ctx context.Context, variant entity.Variant) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (entity.Variant, error)
}

func New(r Repository) Service {
	return Service{repo: r}
}

func (s Service) Create(ctx context.Context, req variantparam.AddVariantRequest) (variantparam.AddVariantResponse, error) {
	const op = "variantService.Create"
	v := entity.Variant{
		Price:    req.Price,
		Quantity: req.Quantity,
		VariantAttribute: []entity.VariantAttribute{entity.VariantAttribute{
			Key:      req.Key,
			Value:    req.Value,
			Label:    req.Label,
			HexColor: req.HexColor,
		}},
	}

	v.BeforeCreate()
	if err := s.repo.Create(ctx, v); err != nil {
		return variantparam.AddVariantResponse{}, richerror.New(op).WithWarpError(err)
	}

	return variantparam.AddVariantResponse{}, nil
}

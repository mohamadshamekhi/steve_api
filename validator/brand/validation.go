package brandvalidator

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
)

type Repository interface {
	CheckBrandNameIsExist(ctx context.Context, brandName string) (bool, error)
	CheckBrandIDIsExist(ctx context.Context, id string) (bool, error)
	FindByName(ctx context.Context, name string) (entity.Brand, error)
}
type Validator struct {
	repo     Repository
	mediaSvc mediaservice.Service
}

func New(repo Repository, mediaSvc mediaservice.Service) Validator {
	return Validator{repo: repo, mediaSvc: mediaSvc}
}

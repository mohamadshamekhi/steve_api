package categoryvalidator

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Category, error)
}
type Validator struct {
	repo     Repository
	mediaSvc mediaservice.Service
}

func New(repo Repository, mediaSvc mediaservice.Service) Validator {
	return Validator{repo: repo, mediaSvc: mediaSvc}
}

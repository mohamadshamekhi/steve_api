package categoryvalidator

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Category, error)
}
type Validator struct {
	repo         Repository
	mediaSvc     mediaservice.Service
	AttrProvider AttrProvider
}

type AttrProvider interface {
	FindByID(ctx context.Context, req attributeparam.GetAttributeRequest) (attributeparam.GetAttributeResponse, error)
}

func New(repo Repository, mediaSvc mediaservice.Service, AttrProvider AttrProvider) Validator {
	return Validator{repo: repo, mediaSvc: mediaSvc, AttrProvider: AttrProvider}
}

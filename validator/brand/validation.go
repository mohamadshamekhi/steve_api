package brandvalidator

import "context"

type Repository interface {
	CheckBrandNameIsExist(ctx context.Context, brandName string) (bool, error)
}
type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

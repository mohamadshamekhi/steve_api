package mediavalidator

import "context"

type Repository interface {
	CheckMediaIsExist(ctx context.Context, id string) (bool, error)
}
type Validator struct {
	repo Repository
}

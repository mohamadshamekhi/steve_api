package commentvalidator

import "golang.org/x/net/context"

type Validator struct {
	repo Repository
}

type Repository interface {
	CommentIsExist(ctx context.Context, id string) (bool, error)
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

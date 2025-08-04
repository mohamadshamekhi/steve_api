package backofficevalidatorauth

import "context"

type Validator struct {
	repo Repository
}

type Repository interface {
	IsBackOfficeUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

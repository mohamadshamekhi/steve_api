package backofficevalidatorauth

import "context"

const (
	phoneNumberRegex = `^\+?[(]?\d{3}[)]?[-\s.]?\d{3}[-\s.]?\d{4,6}$`
)

type Validator struct {
	repo Repository
}

type Repository interface {
	IsBackOfficeUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

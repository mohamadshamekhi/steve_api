package backofficevalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Validator struct {
	repo Repository
}

type Repository interface {
	IsBackOfficeUserIsExistByID(ctx context.Context, id string) (bool, error)
	IsBackOfficeUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	GetBackOfficeUserByID(ctx context.Context, id string) (entity.BackOffice, error)
	GetBackOfficeUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.BackOffice, error)
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

func (v Validator) isBackOfficeUserIsExist(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return richerror.New("validator.isBackOfficeUserIsExist").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type assertion failed"))
	}
	isExist, err := v.repo.IsBackOfficeUserIsExistByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("validation.phone_number_is_already_exist")
	}

	return nil
}

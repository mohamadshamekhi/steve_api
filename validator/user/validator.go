package uservalidator

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validator struct {
	repo Repository
}

type Repository interface {
	CheckUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	CheckUserIsExistByID(ctx context.Context, id string) (bool, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error)
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

func (v Validator) checkPhoneNumberIsUniqueForUpdate(userID string) validation.RuleFunc {
	return func(value interface{}) error {
		var phonenumber string
		switch v := value.(type) {
		case string:
			phonenumber = v
		case *string:
			if v == nil {
				return nil
			}
			phonenumber = *v
		default:
			return richerror.New("uservalidator.type_assertion").
				WithKind(richerror.KindUnexpected).
				WithWarpError(errors.New("unexpected type for brand name"))
		}

		existingUser, err := v.repo.GetUserByPhoneNumber(context.Background(), phonenumber)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return nil
			}
			return err
		}

		if existingUser.ID != userID {
			return errors.New("validation.phone_number_is_already_exist")
		}

		return nil
	}
}

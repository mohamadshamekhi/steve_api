package attributedefinitionvalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Repository interface {
	FindByKey(ctx context.Context, key string) (entity.AttributeDefinition, error)
	FindByID(ctx context.Context, id string) (entity.AttributeDefinition, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) checkKeyIsExist(value any) error {
	var key string
	key, _ = value.(string)
	_, err := v.repo.FindByKey(context.Background(), key)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}
	return nil
}
func (v Validator) checkKeyIsUnique(value any) error {
	var key string
	key, _ = value.(string)
	_, err := v.repo.FindByKey(context.Background(), key)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return nil
			}
		}
		return err
	}
	return errors.New("validation.attr_def_key_is_already_exist")
}

func (v Validator) checkAttrIsExistByID(value any) error {
	var id string
	id, _ = value.(string)
	_, err := v.repo.FindByID(context.Background(), id)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) {
			if rErr.Kind() == richerror.KindNotFound {
				return errors.New(rErr.MessageKey())
			}
		}
		return err
	}
	return nil
}

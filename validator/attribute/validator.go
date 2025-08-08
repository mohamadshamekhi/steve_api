package attributevalidator

import (
	"context"
	"errors"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.AttributeSet, error)
}
type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) checkAttrIsExist(value any) error {
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

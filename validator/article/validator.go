package articlevalidator

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Article, error)
}

type Validator struct {
	repo     Repository
	mediaSvc mediaservice.Service
}

func New(repo Repository, mediaSvc mediaservice.Service) Validator {
	return Validator{repo: repo, mediaSvc: mediaSvc}
}

func (v Validator) articleIsExist(value any) error {
	var id string
	switch v := value.(type) {
	case string:
		id = v
	case *string:
		if v == nil {
			return nil
		}
		id = *v
	default:
		return richerror.New("articlevalidator.articleIsExist.type_assertion").
			WithKind(richerror.KindUnexpected).
			WithWarpError(errors.New("unexpected type for brand name"))
	}

	if id == "" {
		return nil
	}

	a, err := v.repo.FindByID(context.Background(), id)

	if err != nil || a.ID == "" {
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

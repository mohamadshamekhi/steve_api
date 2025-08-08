package mongodbarticle

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, article entity.Article) error {
	const op = "mongodbarticle.Create"
	_, err := d.collection.InsertOne(ctx, article)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]any{"article": article}).WithKind(richerror.KindUnexpected)
	}

	return nil
}

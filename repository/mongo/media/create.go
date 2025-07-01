package mongodbmedia

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) AddMedia(ctx context.Context, media entity.Media) (entity.Media, error) {
	const op = "mongodb.addMedia"

	if _, err := d.collection.InsertOne(ctx, media); err != nil {
		return entity.Media{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	return media, nil
}

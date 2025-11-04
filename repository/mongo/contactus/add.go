package mongodbcontactus

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, cm entity.ContactUs) error {
	const op = "mongodbcontactus.Create"
	_, err := d.collection.InsertOne(ctx, cm)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"contactus": cm}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

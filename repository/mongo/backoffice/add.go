package mongodbbackoffice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) AddNewBackofficeUser(ctx context.Context, u entity.BackOffice) error {
	const op = "mongodbbackoffice.AddNewBackofficeUser"
	_, err := d.collection.InsertOne(ctx, u)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithMessageKey("global.something_wrong").
			WithMeta(map[string]any{"backoffice": u}).WithKind(richerror.KindUnexpected)
	}
	return nil
}

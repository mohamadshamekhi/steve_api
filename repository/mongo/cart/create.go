package mongodbcart

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) CreateCart(ctx context.Context, cart entity.Cart) error {
	const op = "mongodbcart.CreateCart"

	_, err := d.collection.InsertOne(ctx, cart)
	if err != nil {
		return richerror.New(op).
			WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}

	return nil
}

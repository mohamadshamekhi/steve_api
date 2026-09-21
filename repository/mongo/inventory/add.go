package mongodbinventory

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (d DB) Create(ctx context.Context, m entity.StockMovement) error {
	const op = "mongodbinventory.Create"

	if _, err := d.collection.InsertOne(ctx, m); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

// CreateMany appends a batch of movements (one purchase invoice, one order)
// in a single round trip.
func (d DB) CreateMany(ctx context.Context, movements []entity.StockMovement) error {
	const op = "mongodbinventory.CreateMany"

	if len(movements) == 0 {
		return nil
	}
	docs := make([]any, 0, len(movements))
	for _, m := range movements {
		docs = append(docs, m)
	}
	if _, err := d.collection.InsertMany(ctx, docs); err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).
			WithMessageKey("global.something_wrong")
	}
	return nil
}

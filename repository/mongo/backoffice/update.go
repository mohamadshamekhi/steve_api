package mongodbbackoffice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) UpdateBackofficeUser(ctx context.Context, backoffice entity.BackOffice) error {
	const op = "mongodbbackoffice.UpdateBackofficeUser"
	filter := bson.M{"_id": backoffice.ID}
	backoffice.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, filter, backoffice)
	if err != nil {
		return richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	return nil
}

package mongodbpurchase

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

func (d DB) Update(ctx context.Context, p entity.Purchase) error {
	const op = "mongodbpurchase.Update"

	p.UpdatedAt = time.Now()
	result, err := d.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("purchase.not_found")
	}
	return nil
}

// UpdateStatus flips the invoice lifecycle without rewriting the whole
// document, so a concurrent item edit can't be silently clobbered.
func (d DB) UpdateStatus(ctx context.Context, id string, status entity.PurchaseStatus, confirmedAt *time.Time) error {
	const op = "mongodbpurchase.UpdateStatus"

	set := bson.M{"status": status, "updated_at": time.Now()}
	if confirmedAt != nil {
		set["confirmed_at"] = *confirmedAt
	}

	result, err := d.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": set})
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("purchase.not_found")
	}
	return nil
}

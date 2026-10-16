package mongodbproduct

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// SetPartnerConfig tunes a single product's partner (colleague-sales) pricing.
//   - eligible nil  → leave the include/exclude flags untouched
//   - eligible true → sell to partners (partner_eligible=true, partner_excluded=false)
//   - eligible false→ exclude from partner sales (partner_eligible=false, partner_excluded=true)
//   - percent nil   → clear the per-product discount override
func (d DB) SetPartnerConfig(ctx context.Context, productID string, eligible *bool, percent *uint) error {
	const op = "mongodbproduct.SetPartnerConfig"

	set := bson.M{"updated_at": time.Now()}
	unset := bson.M{}

	if eligible != nil {
		set["partner_eligible"] = *eligible
		set["partner_excluded"] = !*eligible
	}
	if percent == nil {
		unset["partner_discount_percent"] = ""
	} else {
		set["partner_discount_percent"] = *percent
	}

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	result, err := d.collection.UpdateOne(ctx, bson.M{"_id": productID}, update)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	if result.MatchedCount == 0 {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}
	return nil
}

// BulkSetPartnerEligible includes (eligible=true) or excludes (eligible=false)
// many products from partner sales at once.
func (d DB) BulkSetPartnerEligible(ctx context.Context, productIDs []string, eligible bool) error {
	const op = "mongodbproduct.BulkSetPartnerEligible"
	_, err := d.collection.UpdateMany(ctx,
		bson.M{"_id": bson.M{"$in": productIDs}},
		bson.M{"$set": bson.M{
			"partner_eligible": eligible,
			"partner_excluded": !eligible,
			"updated_at":       time.Now(),
		}},
	)
	if err != nil {
		return richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return nil
}

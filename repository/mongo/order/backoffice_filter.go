package mongodborder

import (
	"context"
	"regexp"
	"strings"

	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/phonenumber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// maxPhoneMatches caps how many users a partial phone search can fan out to.
const maxPhoneMatches = 500

// backofficeOrderFilter turns the admin list filter into a Mongo filter.
// empty is true when the filter can match nothing (a phone with no user), so
// the caller can skip the query.
func (d DB) backofficeOrderFilter(ctx context.Context, f orderparam.BackofficeFilterOrder) (bson.M, bool, error) {
	filter := bson.M{}

	if f.UserID != nil && *f.UserID != "" {
		filter["user_id"] = *f.UserID
	}
	if f.OrderID != nil && *f.OrderID != "" {
		filter["_id"] = strings.TrimSpace(*f.OrderID)
	}
	if f.Status != nil && *f.Status != "" {
		statuses := splitCSV(*f.Status)
		if len(statuses) == 1 {
			filter["status"] = statuses[0]
		} else if len(statuses) > 1 {
			filter["status"] = bson.M{"$in": statuses}
		}
	}
	if f.PurchaseType != nil && *f.PurchaseType != "" {
		filter["purchase_type"] = *f.PurchaseType
	}
	if f.IsPartnerOrder != nil {
		if *f.IsPartnerOrder {
			filter["is_partner_order"] = true
		} else {
			filter["is_partner_order"] = bson.M{"$ne": true}
		}
	}

	created := bson.M{}
	if f.FromTime != nil {
		created["$gte"] = *f.FromTime
	}
	if f.ToTime != nil {
		created["$lt"] = *f.ToTime
	}
	if len(created) > 0 {
		filter["created_at"] = created
	}

	amount := bson.M{}
	if f.MinAmount != nil {
		amount["$gte"] = *f.MinAmount
	}
	if f.MaxAmount != nil {
		amount["$lte"] = *f.MaxAmount
	}
	if len(amount) > 0 {
		filter["final_amount"] = amount
	}

	if f.Phone != nil && strings.TrimSpace(*f.Phone) != "" {
		phone := strings.TrimSpace(*f.Phone)
		// A full number in any accepted shape (+98…, 0098…) is normalised to
		// the stored form; a partial one is matched as typed.
		if n := phonenumber.Normalize(phone); n != "" {
			phone = n
		}
		ids, err := d.userIDsByPhone(ctx, phone)
		if err != nil {
			return nil, false, err
		}
		if len(ids) == 0 {
			return nil, true, nil
		}
		if uid, ok := filter["user_id"].(string); ok {
			if !contains(ids, uid) {
				return nil, true, nil
			}
		} else {
			filter["user_id"] = bson.M{"$in": ids}
		}
	}

	return filter, false, nil
}

// userIDsByPhone reads the users collection directly: orders only store
// user_id, and a join through the user service would page through all users.
func (d DB) userIDsByPhone(ctx context.Context, phone string) ([]string, error) {
	users := d.collection.Database().Collection("users")
	opts := options.Find().SetProjection(bson.M{"_id": 1}).SetLimit(maxPhoneMatches)
	cursor, err := users.Find(ctx, bson.M{"phone_number": bson.M{"$regex": regexp.QuoteMeta(phone)}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

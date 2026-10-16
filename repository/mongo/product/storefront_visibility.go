package mongodbproduct

import (
	"github.com/abolfazlnorzad/steve/entity"
	"go.mongodb.org/mongo-driver/bson"
)

// websiteChannelCondition matches products that may appear on the storefront.
//
// A product written before sales channels existed has no `sales_channels`
// field (or an empty one) and was, by definition, a website product — so both
// of those shapes must still match. Only an explicit channel list that leaves
// "website" out hides the product, which is how an in-shop-only item and a
// service sold over the counter stay off the site.
func websiteChannelCondition() bson.M {
	return bson.M{"$or": []bson.M{
		{"sales_channels": bson.M{"$exists": false}},
		{"sales_channels": bson.M{"$size": 0}},
		{"sales_channels": entity.SalesChannelWebsite},
	}}
}

// withWebsiteChannel adds the storefront visibility gate to a bson.M filter
// without clobbering an `$or` the caller already set, by folding both into
// `$and`.
func withWebsiteChannel(filter bson.M) bson.M {
	conditions := []bson.M{websiteChannelCondition()}

	if existing, ok := filter["$and"]; ok {
		switch typed := existing.(type) {
		case []bson.M:
			conditions = append(conditions, typed...)
		case bson.A:
			for _, c := range typed {
				if m, isMap := c.(bson.M); isMap {
					conditions = append(conditions, m)
				}
			}
		}
		delete(filter, "$and")
	}
	if existing, ok := filter["$or"]; ok {
		conditions = append(conditions, bson.M{"$or": existing})
		delete(filter, "$or")
	}

	filter["$and"] = conditions
	return filter
}

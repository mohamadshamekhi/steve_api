package mongodbuser

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"go.mongodb.org/mongo-driver/bson"
)

// FindUsersByBirthMonthDay returns users whose birth_date falls on any of the
// given month/day pairs, regardless of birth year. birth_date is stored as UTC
// midnight, so $month/$dayOfMonth (UTC) match the calendar date entered.
func (d DB) FindUsersByBirthMonthDay(ctx context.Context, days []entity.MonthDay) ([]entity.User, error) {
	const op = "mongodbuser.FindUsersByBirthMonthDay"
	if len(days) == 0 {
		return nil, nil
	}

	or := make(bson.A, 0, len(days))
	for _, md := range days {
		or = append(or, bson.M{"$and": bson.A{
			bson.M{"$eq": bson.A{bson.M{"$month": "$birth_date"}, int(md.Month)}},
			bson.M{"$eq": bson.A{bson.M{"$dayOfMonth": "$birth_date"}, md.Day}},
		}})
	}
	filter := bson.M{
		"birth_date": bson.M{"$type": "date"},
		"$expr":      bson.M{"$or": or},
	}

	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	defer cursor.Close(ctx)

	var users []entity.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, richerror.New(op).WithWarpError(err).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	return users, nil
}

package entity

import "time"

type Rating struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Rating    float64   `bson:"rating"`
	CreatedAt time.Time `bson:"created_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

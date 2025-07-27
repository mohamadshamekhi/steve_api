package entity

type User struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	PhoneNumber string `bson:"phone_number"`
}

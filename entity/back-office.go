package entity

type BackOffice struct {
	ID           string `bson:"_id"`
	Name         string `bson:"name"`
	PhoneNumber  string `bson:"phone_number"`
	NationalCode string `bson:"national_code"`
	Address      string `bson:"address"`
	ZipCode      string `bson:"zip_code"`
	AvatarID     int64  `bson:"avatar_id"`
}

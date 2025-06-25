package entity

type Brand struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Description string `bson:"description,omitempty"`
	MediaID     string `bson:"media_id,omitempty"`
}

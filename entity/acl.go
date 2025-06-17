package entity

type Permission struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

type Role struct {
	ID            string   `bson:"_id"`
	Name          string   `bson:"name"`
	Description   string   `bson:"description"`
	PermissionIDs []string `bson:"permission_ids"`
}

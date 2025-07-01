package entity

import "github.com/oklog/ulid/v2"

type Permission struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

func (p *Permission) BeforeCreate() {
	id := ulid.Make()
	p.ID = id.String()
}

type Role struct {
	ID             string   `bson:"_id"`
	Name           string   `bson:"name"`
	Description    string   `bson:"description"`
	PermissionKeys []string `bson:"permission_keys"` // بهترین روش برای این کار
}

func (r *Role) BeforeCreate() {
	id := ulid.Make()
	r.ID = id.String()
}

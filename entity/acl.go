package entity

import "github.com/oklog/ulid/v2"

type Permission struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

func (p *Permission) BeforeCreate() {
	id := ulid.Make()
	p.ID = id.String()
}

type Role struct {
	ID             string   `bson:"_id" json:"id"`
	Name           string   `bson:"name" json:"name"`
	Description    string   `bson:"description" json:"description"`
	PermissionKeys []string `bson:"permission_keys" json:"permission_keys"`
}

func (r *Role) BeforeCreate() {
	id := ulid.Make()
	r.ID = id.String()
}

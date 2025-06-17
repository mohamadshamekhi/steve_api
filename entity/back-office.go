package entity

type BackOffice struct {
	ID            string   `bson:"_id"`
	Name          string   `bson:"name"`
	PhoneNumber   string   `bson:"phone_number"`
	NationalCode  string   `bson:"national_code"`
	Address       string   `bson:"address"`
	ZipCode       string   `bson:"zip_code"`
	AvatarID      int64    `bson:"avatar_id"`
	IsActive      bool     `bson:"is_active"`
	isSuperAdmin  bool     `bson:"is_super_admin"`
	RoleIDs       []string `bson:"role_ids,omitempty"`
	PermissionIDs []string `bson:"permission_ids,omitempty"`
}

func (b BackOffice) IsSuperAdmin() bool {
	return b.isSuperAdmin
}

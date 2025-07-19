package backofficeparam

type UpdateBackofficeUserRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	NationCode  string   `json:"nation_code,omitempty"`
	Address     string   `json:"address"`
	ZipCode     string   `json:"zip_code,omitempty"`
	IsActive    *bool    `json:"is_active"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateBackofficeUserResponse struct {
}

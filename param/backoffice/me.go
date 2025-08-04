package backofficeparam

type MeRequest struct {
	ID string `json:"id"`
}

type MeResponse struct {
	Backoffice BackofficeInfo `json:"backoffice"`
}
type BackofficeInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PhoneNumber  string   `json:"phone_number"`
	NationalCode string   `json:"national_code"`
	Address      string   `json:"address"`
	ZipCode      string   `json:"zip_code"`
	IsActive     bool     `json:"is_active"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
}

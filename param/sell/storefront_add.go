package sellparam

type AddSellRequest struct {
	UserID            string `json:"user_id"`
	BrandID           string `json:"brand_id"`
	ProductID         string `json:"product_id"`
	Color             string `json:"color"`
	BatteryHealth     string `json:"battery_health"`
	RegisterStatus    string `json:"register_status"`
	PhysicalCondition string `json:"physical_condition"`
	Description       string `json:"description"`
}

type AddSellResponse struct {
}

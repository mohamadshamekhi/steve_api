package entity

import "time"

type SellStatus string

type Sell struct {
	ID                string     `bson:"_id"`
	UserID            string     `bson:"user_id"`
	BrandID           string     `bson:"brand_id"`
	ProductID         string     `bson:"product_id"`
	Color             string     `bson:"color"`
	BatteryHealth     string     `bson:"battery_health"`
	RegisterStatus    string     `bson:"register_status"`
	PhysicalCondition string     `bson:"physical_condition"`
	Description       string     `bson:"description"`
	Status            SellStatus `bson:"status"`
	CreatedAt         time.Time  `bson:"created_at"`
	UpdatedAt         time.Time  `bson:"updated_at"`
}

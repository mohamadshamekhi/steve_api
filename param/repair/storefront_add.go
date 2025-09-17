package repairparam

import "github.com/abolfazlnorzad/steve/entity"

type AddRepairRequest struct {
	UserID      string                 `json:"user_id"`
	BrandID     string                 `json:"brand_id"`
	ProductID   string                 `json:"product_id"`
	Description string                 `json:"description"`
	Galleries   []entity.RepairGallery `bson:"galleries"`
}

type AddRepairResponse struct {
}

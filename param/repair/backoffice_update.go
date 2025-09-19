package repairparam

import "github.com/abolfazlnorzad/steve/entity"

type UpdateRepairRequest struct {
	ID     string              `json:"id"`
	Status entity.RepairStatus `json:"status"`
}

type UpdateRepairResponse struct {
}

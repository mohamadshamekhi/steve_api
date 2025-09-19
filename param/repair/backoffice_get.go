package repairparam

type GetRepairRequest struct {
	ID string
}

type GetRepairResponse struct {
	RepairInfo RepairInfo `json:"repair"`
}

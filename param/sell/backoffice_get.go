package sellparam

type GetSellRequest struct {
	ID string
}

type GetSellResponse struct {
	Sell SellInfo `json:"sell"`
}

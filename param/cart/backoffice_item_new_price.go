package cartparam

// SetItemUnitPriceBackofficeRequest adjusts a single invoice line by storing
// the unit price the salesperson wants to charge. The backend derives the
// manual_discount (positive = discount, negative = markup) from the diff with
// the original price_at_time.
type SetItemUnitPriceBackofficeRequest struct {
	UserID    string `json:"user_id"`
	VariantID string `json:"variant_id"`
	UnitPrice int64  `json:"unit_price"`
}

type SetItemUnitPriceBackofficeResponse struct{}

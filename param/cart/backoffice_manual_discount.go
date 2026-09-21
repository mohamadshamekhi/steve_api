package cartparam

// SetItemManualDiscountBackofficeRequest sets the discount given by hand on a
// single invoice line, in Toman for the whole line. Zero clears it.
type SetItemManualDiscountBackofficeRequest struct {
	UserID    string `json:"user_id"`
	VariantID string `json:"variant_id"`
	Amount    uint   `json:"amount"`
}

type SetItemManualDiscountBackofficeResponse struct{}

// SetCartManualDiscountBackofficeRequest sets the discount given by hand on
// the invoice as a whole, applied after the per-line ones. Zero clears it.
type SetCartManualDiscountBackofficeRequest struct {
	UserID string `json:"user_id"`
	Amount uint   `json:"amount"`
}

type SetCartManualDiscountBackofficeResponse struct{}

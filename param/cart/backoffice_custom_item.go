package cartparam

import "github.com/abolfazlnorzad/steve/entity"

// UpsertCustomItemBackofficeRequest adds or edits a hand-typed line on an
// in-person invoice: a repair fee, an Apple ID, or a part bought from outside
// that the catalog does not carry. Leave VariantID empty to add a new line;
// pass the id returned earlier to edit that line (price, title, quantity).
type UpsertCustomItemBackofficeRequest struct {
	UserID string `json:"user_id"`
	// VariantID is the line id returned when the line was added. Empty means
	// "create a new line".
	VariantID string `json:"variant_id"`
	// Title is what gets printed on the invoice, e.g. «هزینه تعمیر آیفون ۱۳».
	Title string `json:"title"`
	// UnitPrice is what the customer pays per unit, in Toman.
	UnitPrice int64 `json:"unit_price"`
	// Quantity is the number of units; the panel's +/- buttons resend this.
	Quantity int `json:"quantity"`
	// Kind separates a service (repair, Apple ID) from a physical part, which
	// decides how the line is booked in the inventory ledger.
	Kind entity.ProductKind `json:"kind"`
	// UnitCost is what the shop paid per unit for a part bought outside; it
	// never reaches the customer and only feeds the profit report.
	UnitCost uint `json:"unit_cost"`
}

type UpsertCustomItemBackofficeResponse struct {
	// VariantID identifies the line for later edits, discounting or removal.
	VariantID string `json:"variant_id"`
}

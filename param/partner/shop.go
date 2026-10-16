package partnerparam

// --- partner shop: cart ---

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity"`
}

type RemoveFromCartRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

type CartItemView struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id"`
	ProductName string `json:"product_name"`
	URL         string `json:"url"`
	Quantity    int    `json:"quantity"`
	BasePrice   uint64 `json:"base_price"`
	UnitPrice   uint64 `json:"unit_price"`
	LineTotal   uint64 `json:"line_total"`
	InStock     bool   `json:"in_stock"`
}

type CartResponse struct {
	CartID     string         `json:"cart_id"`
	Items      []CartItemView `json:"items"`
	Subtotal   uint64         `json:"subtotal"`
	SavedTotal uint64         `json:"saved_total"`
}

// --- partner shop: checkout & payment ---

const (
	PaymentMethodOnline   = "online"
	PaymentMethodInPerson = "in_person"
)

type CheckoutRequest struct {
	PaymentMethod string `json:"payment_method"` // "online" | "in_person"
}

type CheckoutResponse struct {
	OrderID              string `json:"order_id"`
	PaymentMethod        string `json:"payment_method"`
	PaymentURL           string `json:"payment_url,omitempty"`
	GatewayTransactionID string `json:"gateway_transaction_id,omitempty"`
}

type VerifyPaymentRequest struct {
	GatewayTransactionID string `json:"gateway_transaction_id"`
}

// --- partner order fulfilment (partner_manager) ---

type RecordInPersonPaymentRequest struct {
	Description string `json:"description"`
}

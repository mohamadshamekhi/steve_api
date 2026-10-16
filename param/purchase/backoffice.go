package purchaseparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// CreatePurchaseItemRequest is one line of a purchase invoice.
type CreatePurchaseItemRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Quantity  uint64 `json:"quantity"`
	UnitCost  uint64 `json:"unit_cost"`
	// AddToStock defaults to true; set it false for goods bought and handed
	// straight to a customer without ever entering the warehouse.
	AddToStock *bool `json:"add_to_stock"`
}

type CreatePurchaseRequest struct {
	SupplierID    string                      `json:"supplier_id"`
	InvoiceNumber string                      `json:"invoice_number"`
	Items         []CreatePurchaseItemRequest `json:"items"`
	PaidAmount    uint64                      `json:"paid_amount"`
	Note          string                      `json:"note"`
	OrderID       string                      `json:"order_id"`
	PurchasedAt   *time.Time                  `json:"purchased_at"`
	// Confirm immediately applies the invoice to stock instead of leaving it
	// as a draft.
	Confirm bool `json:"confirm"`
}

type CreatePurchaseResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

type UpdatePurchaseRequest struct {
	ID            string                      `json:"-"`
	SupplierID    *string                     `json:"supplier_id"`
	InvoiceNumber *string                     `json:"invoice_number"`
	Items         []CreatePurchaseItemRequest `json:"items"`
	PaidAmount    *uint64                     `json:"paid_amount"`
	Note          *string                     `json:"note"`
	PurchasedAt   *time.Time                  `json:"purchased_at"`
}

type UpdatePurchaseResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

type ListPurchasesRequest struct {
	SupplierID   string
	Status       entity.PurchaseStatus
	PaymentState entity.PaymentState
	OrderID      string
	From         *time.Time
	To           *time.Time
	Pagination   commonparam.PaginationRequest
}

type ListPurchasesResponse struct {
	Purchases  []entity.Purchase              `json:"purchases"`
	Pagination commonparam.PaginationResponse `json:"pagination"`
}

type GetPurchaseResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

type ConfirmPurchaseResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

type CancelPurchaseResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

// SetPaidAmountRequest records how much of the invoice the shop has paid.
type SetPaidAmountRequest struct {
	// PaymentMethod is how the shop settled it, recorded in the cashbox.
	PaymentMethod entity.PaymentMethod `json:"payment_method"`
	ID            string               `json:"-"`
	PaidAmount    uint64               `json:"paid_amount"`
}

type SetPaidAmountResponse struct {
	Purchase entity.Purchase `json:"purchase"`
}

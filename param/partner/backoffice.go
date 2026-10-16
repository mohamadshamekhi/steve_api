package partnerparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type ListRequest struct {
	Pagination commonparam.PaginationRequest `json:"pagination" query:"pagination"`
	Status     string                        `json:"status" query:"status"`
}

type ListResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Partners   []entity.Partner               `json:"partners"`
}

type GetResponse struct {
	Partner entity.Partner `json:"partner"`
}

type ApproveRequest struct {
	DiscountPercent uint `json:"discount_percent"`
}

type RejectRequest struct {
	Reason string `json:"reason"`
}

type UpdatePricingRequest struct {
	DiscountPercent uint `json:"discount_percent"`
}

type EmptyResponse struct{}

// SetProductConfigRequest tunes a single product's partner pricing.
//   - eligible = true  → sell this product to partners
//   - eligible = false → exclude this product from partner sales
//   - eligible omitted → leave include/exclude unchanged (only touch percent)
type SetProductConfigRequest struct {
	Eligible        *bool `json:"eligible"`
	DiscountPercent *uint `json:"discount_percent"`
}

// BulkSetEligibleRequest includes or excludes many products at once.
type BulkSetEligibleRequest struct {
	ProductIDs []string `json:"product_ids"`
	Eligible   bool     `json:"eligible"`
}

// --- admin: product partner-config listing ---

type ProductListRequest struct {
	Pagination commonparam.PaginationRequest
	Name       string
	BrandID    string
	CategoryID string
	Eligible   *bool
}

type ProductListVariant struct {
	ID       string `json:"id"`
	Price    uint64 `json:"price"`
	Quantity uint64 `json:"quantity"`
}

type ProductListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	IsPublished bool   `json:"is_published"`
	// SoldToPartners is the effective state: is this product currently
	// purchasable at the partner price (after program mode + exclusion)?
	SoldToPartners bool `json:"sold_to_partners"`
	// PartnerExcluded is the raw admin exclusion flag.
	PartnerExcluded        bool                 `json:"partner_excluded"`
	PartnerDiscountPercent *uint                `json:"partner_discount_percent"`
	Variants               []ProductListVariant `json:"variants"`
}

type ProductListResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Products   []ProductListItem              `json:"products"`
}

// --- admin: partner order listing ---

type OrderListRequest struct {
	Pagination commonparam.PaginationRequest
	Status     string
}

type OrderListItem struct {
	ID               string                   `json:"id"`
	PartnerID        string                   `json:"partner_id"`
	UserID           string                   `json:"user_id"`
	Status           string                   `json:"status"`
	PurchaseType     string                   `json:"purchase_type"`
	SubtotalAmount   uint                     `json:"subtotal_amount"`
	FinalAmount      uint                     `json:"final_amount"`
	AppliedDiscounts []entity.AppliedDiscount `json:"applied_discounts"`
	Items            []entity.OrderItem       `json:"items"`
	HandoverBy       string                   `json:"handover_by"`
	HandoverAt       *time.Time               `json:"handover_at"`
	CreatedAt        time.Time                `json:"created_at"`
}

type OrderListResponse struct {
	Pagination commonparam.PaginationResponse `json:"pagination"`
	Orders     []OrderListItem                `json:"orders"`
}

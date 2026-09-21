package productparam

import "github.com/abolfazlnorzad/steve/entity"

type CreateProductRequest struct {
	Name         string `json:"name"`
	ENTitle      string `json:"en_title"`
	BrandID      string `json:"brand_id"`
	CategoryID   string `json:"category_id"`
	IsSecondHand *bool  `json:"is_second_hand"`
	Description  string `json:"description"`
	// Kind marks the item as a stockable good or a service the shop performs
	// (repair, Apple-ID, ...). Defaults to physical.
	Kind entity.ProductKind `json:"kind"`
	// SalesChannels controls where it may be sold. Defaults to both website
	// and in-person; pass only ["in_person"] for an item the shop sells over
	// the counter but does not list online.
	SalesChannels []entity.SalesChannel `json:"sales_channels"`
	// Price is the service's fee. Required when Kind is service and ignored
	// for goods, whose prices live on their variants. A service is created
	// complete and sellable in one call — the caller never deals with the
	// variant the price is stored on.
	Price uint64 `json:"price"`
}
type CreateProductResponse struct {
	Product ProductCreated `json:"product"`
}

// ProductCreated is what the create/update endpoints echo back. It is a
// dedicated DTO rather than the raw entity so the keys stay snake_case like
// every other product response — the entity carries only bson tags, which
// would serialise as PascalCase and break the caller's expectations.
type ProductCreated struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	ENTitle       string                `json:"en_title"`
	Slug          string                `json:"slug"`
	BrandID       string                `json:"brand_id,omitempty"`
	CategoryID    string                `json:"category_id,omitempty"`
	Kind          entity.ProductKind    `json:"kind"`
	SalesChannels []entity.SalesChannel `json:"sales_channels"`
	IsPublished   bool                  `json:"is_published"`
	IsSecondHand  bool                  `json:"is_second_hand"`
	Description   string                `json:"description"`
	// Price is echoed back for a service. Omitted for goods, which price per
	// variant instead.
	Price uint64 `json:"price,omitempty"`
	// VariantID is the service's implicit pricing row, exposed so an old
	// client that still speaks variants can find it. New callers should use
	// Price and the product-level update instead.
	VariantID string `json:"variant_id,omitempty"`
}

// NewProductCreated normalises the entity's defaulted fields (kind and sales
// channels) so the caller never has to re-derive them.
func NewProductCreated(p entity.Product) ProductCreated {
	created := ProductCreated{
		ID:            p.ID,
		Name:          p.Name,
		ENTitle:       p.ENTitle,
		Slug:          p.Slug,
		BrandID:       p.BrandID,
		CategoryID:    p.CategoryID,
		Kind:          p.GetKind(),
		SalesChannels: p.GetSalesChannels(),
		IsPublished:   p.IsPublished,
		IsSecondHand:  p.IsSecondHand,
		Description:   p.Description,
	}
	if v, ok := p.ServiceVariant(); ok {
		created.Price = v.Price
		created.VariantID = v.ID
	}
	return created
}

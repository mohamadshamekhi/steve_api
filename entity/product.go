package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// ProductKind separates a real, stockable good from a service the shop
// performs (repair, Apple-ID creation, ...). Services never touch inventory.
type ProductKind string

const (
	ProductKindPhysical ProductKind = "physical" // کالا
	ProductKindService  ProductKind = "service"  // خدمات (تعمیر، ساخت اپل‌آیدی، ...)
)

func (k ProductKind) IsValid() bool {
	return k == ProductKindPhysical || k == ProductKindService
}

// SalesChannel says where an item may be sold. A product that is only
// `in_person` is invisible on the storefront but still usable in backoffice
// invoices and in the shop's accounting.
type SalesChannel string

const (
	SalesChannelWebsite  SalesChannel = "website"   // فروش در وب‌سایت
	SalesChannelInPerson SalesChannel = "in_person" // فقط فروش حضوری / فاکتور
)

func (c SalesChannel) IsValid() bool {
	return c == SalesChannelWebsite || c == SalesChannelInPerson
}

type Product struct {
	ID           string           `bson:"_id"`
	Name         string           `bson:"name"`
	ENTitle      string           `bson:"en_title"`
	Slug         string           `bson:"slug"`
	BrandID      string           `bson:"brand_id"`
	CategoryID   string           `bson:"category_id"`
	Gallery      []ProductGallery `bson:"gallery"`
	IsPublished  bool             `bson:"is_published"`
	IsSecondHand bool             `bson:"is_second_hand"`
	// RequiresFreightShipping marks an oversized product that can only be
	// delivered via freight (باربری), not regular courier.
	RequiresFreightShipping bool            `bson:"requires_freight_shipping,omitempty"`
	Description             string          `bson:"description"`
	Rating                  float64         `bson:"rating"`
	Variants                []Variant       `bson:"variants"`
	Specifications          []Specification `bson:"specifications,omitempty"`
	SalesCount              uint            `bson:"sales_count"`
	// Kind distinguishes a stockable good from a service. Empty value is read
	// as ProductKindPhysical so pre-existing documents keep working.
	Kind ProductKind `bson:"kind,omitempty"`
	// SalesChannels lists where this item may be sold. Empty is read as
	// [website, in_person] for backward compatibility with old documents.
	SalesChannels []SalesChannel `bson:"sales_channels,omitempty"`
	// PartnerEligible explicitly opts the product into partner (B2B) pricing.
	// Decisive only in opt-in mode (partner.all_products_eligible = false).
	PartnerEligible bool `bson:"partner_eligible,omitempty"`
	// PartnerExcluded explicitly opts the product OUT of partner pricing, even
	// in all-products mode. Admin's per-product exclusion.
	PartnerExcluded bool `bson:"partner_excluded,omitempty"`
	// PartnerDiscountPercent overrides the partner's global discount percent
	// for this product. nil → fall back to Partner.DiscountPercent.
	PartnerDiscountPercent *uint     `bson:"partner_discount_percent,omitempty"`
	CreatedAt              time.Time `bson:"created_at"`
	UpdatedAt              time.Time `bson:"updated_at"`
	DeletedAt              time.Time `bson:"deleted_at,omitempty"`
	// PriceReviewedAt is the last time an admin confirmed or changed this
	// product's price. Used by the periodic price-review feature to flag
	// products whose price hasn't been checked in a configured interval.
	PriceReviewedAt time.Time `bson:"price_reviewed_at,omitempty"`
	// SEO fields, editable by the content/SEO agent. All optional so old
	// documents keep working unchanged.
	SeoTitle       string       `bson:"seo_title,omitempty"`
	SeoDescription string       `bson:"seo_description,omitempty"`
	FocusKeyword   string       `bson:"focus_keyword,omitempty"`
	FAQ            []ProductFAQ `bson:"faq,omitempty"`
}
type ProductGallery struct {
	MediaID   string `bson:"media_id" json:"media_id"`
	URL       string `bson:"url" json:"url"`
	IsPrimary bool   `bson:"is_primary" json:"is_primary"`
	Alt       string `bson:"alt,omitempty" json:"alt,omitempty"`
}
type ProductFAQ struct {
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}
type Specification struct {
	Key   string `bson:"key" json:"key"`     // ex: "cpu_processor"
	Value string `bson:"value" json:"value"` // ex: "A18 Bionic"
}

func (p *Product) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.IsPublished = false
	if p.Kind == "" {
		p.Kind = ProductKindPhysical
	}
	if len(p.SalesChannels) == 0 {
		p.SalesChannels = []SalesChannel{SalesChannelWebsite, SalesChannelInPerson}
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	p.PriceReviewedAt = now

	for i := range p.Variants {
		p.Variants[i].ID = ulid.Make().String()
	}
}

func (p Product) GetTitle() string {
	return p.Name
}
func (p Product) GetID() string {
	return p.ID
}

func (p Product) GetMainGallery() ProductGallery {
	if len(p.Gallery) == 0 {
		return ProductGallery{}
	}
	for _, gallery := range p.Gallery {
		if gallery.IsPrimary {
			return gallery
		}
	}
	return p.Gallery[0]
}

// GetKind normalises the zero value of documents written before the field
// existed.
func (p Product) GetKind() ProductKind {
	if p.Kind == "" {
		return ProductKindPhysical
	}
	return p.Kind
}

// IsService reports whether this item is a service (no stock, no shipping).
func (p Product) IsService() bool {
	return p.GetKind() == ProductKindService
}

// GetSalesChannels normalises the zero value of documents written before the
// field existed: they were website+in-person products.
func (p Product) GetSalesChannels() []SalesChannel {
	if len(p.SalesChannels) == 0 {
		return []SalesChannel{SalesChannelWebsite, SalesChannelInPerson}
	}
	return p.SalesChannels
}

// HasChannel reports whether the product may be sold through ch.
func (p Product) HasChannel(ch SalesChannel) bool {
	for _, c := range p.GetSalesChannels() {
		if c == ch {
			return true
		}
	}
	return false
}

// IsSellableOnWebsite is the storefront gate: published, not deleted and
// explicitly opted into the website channel.
func (p Product) IsSellableOnWebsite() bool {
	return p.IsPublished && p.DeletedAt.IsZero() && p.HasChannel(SalesChannelWebsite)
}

// IsSellableInPerson is the backoffice-invoice gate. It deliberately ignores
// IsPublished: an item the shop sells over the counter needs no storefront
// page.
func (p Product) IsSellableInPerson() bool {
	return p.DeletedAt.IsZero() && p.HasChannel(SalesChannelInPerson)
}

// ServiceVariant returns the single pricing row a service is stored with.
//
// A service is modelled as a Product with exactly one variant so that carts,
// orders, invoices and the accounting ledger — all of which address a line as
// product + variant — keep working unchanged. That variant is an
// implementation detail: callers should reach a service's price through here
// rather than indexing into Variants.
func (p Product) ServiceVariant() (Variant, bool) {
	if !p.IsService() || len(p.Variants) == 0 {
		return Variant{}, false
	}
	return p.Variants[0], true
}

// ServicePrice is the price of a service, or zero for a good.
func (p Product) ServicePrice() uint64 {
	v, ok := p.ServiceVariant()
	if !ok {
		return 0
	}
	return v.Price
}

// DisplayPrice is the single price shown for a product in admin lists: the
// service fee for a service, or the lowest variant price for a good (0 when
// it has none yet).
func (p Product) DisplayPrice() uint64 {
	if p.IsService() {
		return p.ServicePrice()
	}
	var min uint64
	for i, v := range p.Variants {
		if i == 0 || v.Price < min {
			min = v.Price
		}
	}
	return min
}

// NewServiceVariant builds the implicit pricing row for a service.
func NewServiceVariant(price uint64) Variant {
	return Variant{
		Price:     price,
		Quantity:  0,
		StockMode: StockModeNone,
	}
}

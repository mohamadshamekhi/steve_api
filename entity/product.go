package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Product struct {
	ID             string           `bson:"_id"`
	Name           string           `bson:"name"`
	ENTitle        string           `bson:"en_title"`
	Slug           string           `bson:"slug"`
	BrandID        string           `bson:"brand_id"`
	CategoryID     string           `bson:"category_id"`
	Gallery        []ProductGallery `bson:"gallery"`
	IsPublished    bool             `bson:"is_published"`
	IsSecondHand   bool             `bson:"is_second_hand"`
	Description    string           `bson:"description"`
	Rating         float64          `bson:"rating"`
	Variants       []Variant        `bson:"variants"`
	Specifications []Specification  `bson:"specifications,omitempty"`
	SalesCount     uint             `bson:"sales_count"`
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
}
type ProductGallery struct {
	MediaID   string `bson:"media_id" json:"media_id"`
	URL       string `bson:"url" json:"url"`
	IsPrimary bool   `bson:"is_primary" json:"is_primary"`
}
type Specification struct {
	Key   string `bson:"key" json:"key"`     // ex: "cpu_processor"
	Value string `bson:"value" json:"value"` // ex: "A18 Bionic"
}

func (p *Product) BeforeCreate() {
	p.ID = ulid.Make().String()
	now := time.Now()
	p.IsPublished = false
	p.CreatedAt = now
	p.UpdatedAt = now

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

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

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
	Description  string           `bson:"description"`

	Variants       []Variant       `bson:"variants"`
	Specifications []Specification `bson:"specifications,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at,omitempty"`
}
type ProductGallery struct {
	MediaID   string `bson:"media_id"`
	URL       string `bson:"url"`
	IsPrimary bool   `bson:"is_primary"`
}
type Specification struct {
	Key   string `bson:"key"`   // ex: "cpu_processor"
	Value string `bson:"value"` // ex: "A18 Bionic"
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

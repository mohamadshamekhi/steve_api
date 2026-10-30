package productparam

import "time"

// CompetitorLink is one row of GET /product/competitor-links: the minimum the
// price agent needs to look a product up on Torob and compare.
type CompetitorLink struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	TorobURL string `json:"torob_url"`
	// OurMinPrice is the cheapest website-sellable variant's list price,
	// before discounts. 0 when the product has no priced variant.
	OurMinPrice uint64 `json:"our_min_price"`
	InStock     bool   `json:"in_stock"`
}

type GetCompetitorLinksRequest struct{}

type GetCompetitorLinksResponse struct {
	Products []CompetitorLink `json:"products"`
}

type UpdateCompetitorPriceRequest struct {
	ProductID   string     `json:"-"`
	MinPrice    uint64     `json:"min_price"`
	MaxPrice    uint64     `json:"max_price"`
	OffersCount uint       `json:"offers_count"`
	CheckedAt   *time.Time `json:"checked_at"`
}

type UpdateCompetitorPriceResponse struct{}

package partnerservice

// Config toggles the colleague-sales (partner) program. Loaded via koanf with
// the ST_PARTNER__ env prefix, e.g. ST_PARTNER__ALLOW_DISCOUNT_STACKING=true.
type Config struct {
	// Enabled turns the whole partner program on/off.
	Enabled bool `koanf:"enabled"`
	// AllowDiscountStacking, when false (default), makes partner checkout
	// ignore public discount codes — the partner price already is the deal.
	AllowDiscountStacking bool `koanf:"allow_discount_stacking"`
	// AllProductsEligible, when true (default), applies the partner discount to
	// every published product; the per-product `partner_eligible` flag then
	// only matters when this is turned off (opt-in mode).
	AllProductsEligible bool `koanf:"all_products_eligible"`
}

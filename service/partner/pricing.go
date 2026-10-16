package partnerservice

// PriceInput carries exactly what the resolver needs to price one variant for a
// partner — nothing more. Callers assemble it from the product, the variant and
// the partner record.
type PriceInput struct {
	BasePrice              uint64
	ProductEligible        bool
	ProductDiscountPercent *uint
	PartnerDiscountPercent uint
}

// PriceResult is the resolver's verdict for one variant.
type PriceResult struct {
	Eligible       bool
	AppliedPercent uint
	FinalPrice     uint64
	DiscountAmount uint64
}

// PartnerPriceResolver computes the price a partner pays for one variant.
//
// Single responsibility: pure pricing arithmetic. It performs no persistence,
// no I/O and holds no state, so it is trivially unit-testable and is the single
// place partner pricing is decided — the partner catalog and (later) partner
// checkout both call this and nothing else.
type PartnerPriceResolver interface {
	Resolve(in PriceInput) PriceResult
}

type percentPriceResolver struct{}

// NewPartnerPriceResolver returns the default percentage-based resolver.
func NewPartnerPriceResolver() PartnerPriceResolver {
	return percentPriceResolver{}
}

func (percentPriceResolver) Resolve(in PriceInput) PriceResult {
	if !in.ProductEligible {
		return PriceResult{
			Eligible:       false,
			AppliedPercent: 0,
			FinalPrice:     in.BasePrice,
			DiscountAmount: 0,
		}
	}

	percent := in.PartnerDiscountPercent
	if in.ProductDiscountPercent != nil {
		percent = *in.ProductDiscountPercent
	}
	if percent > 100 {
		percent = 100
	}

	discount := in.BasePrice * uint64(percent) / 100

	return PriceResult{
		Eligible:       true,
		AppliedPercent: percent,
		FinalPrice:     in.BasePrice - discount,
		DiscountAmount: discount,
	}
}

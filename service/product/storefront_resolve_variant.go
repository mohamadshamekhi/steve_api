package productservice

import (
	"context"
	"strings"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// ResolveVariant backs the printed variant QR label. The variant ID is the
// label's code because it never changes when the product is renamed or its
// slug is regenerated. Unpublished/in-shop-only products still resolve so a
// label on a shelf item never dead-ends; OnWebsite says whether to redirect.
func (s StorefrontService) ResolveVariant(ctx context.Context, req productparam.ResolveVariantRequest) (productparam.ResolveVariantResponse, error) {
	const op = "productservice.ResolveVariant"

	p, err := s.repo.FindByVariantID(ctx, req.VariantID)
	if err != nil {
		return productparam.ResolveVariantResponse{}, richerror.New(op).WithWarpError(err)
	}

	for _, v := range p.Variants {
		if v.ID != req.VariantID {
			continue
		}
		labels := make([]string, 0, len(v.VariantAttribute))
		for _, a := range v.VariantAttribute {
			if a.Label != "" {
				labels = append(labels, a.Label)
			} else if a.Value != "" {
				labels = append(labels, a.Value)
			}
		}
		return productparam.ResolveVariantResponse{
			ProductID:    p.ID,
			VariantID:    v.ID,
			Slug:         p.Slug,
			ProductName:  p.Name,
			VariantLabel: strings.Join(labels, " / "),
			Kind:         string(p.GetKind()),
			OnWebsite:    p.IsSellableOnWebsite(),
		}, nil
	}

	return productparam.ResolveVariantResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
}

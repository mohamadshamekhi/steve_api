package productservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) GetCompetitorLinks(ctx context.Context, _ productparam.GetCompetitorLinksRequest) (productparam.GetCompetitorLinksResponse, error) {
	const op = "productservice.GetCompetitorLinks"

	products, err := s.repo.FindWithTorobURL(ctx)
	if err != nil {
		return productparam.GetCompetitorLinksResponse{}, richerror.New(op).WithWarpError(err)
	}

	links := make([]productparam.CompetitorLink, 0, len(products))
	for _, p := range products {
		minPrice, inStock := ourPriceAndStock(p)
		links = append(links, productparam.CompetitorLink{
			ID:          p.ID,
			Slug:        p.Slug,
			Name:        p.Name,
			TorobURL:    p.TorobURL,
			OurMinPrice: minPrice,
			InStock:     inStock,
		})
	}
	return productparam.GetCompetitorLinksResponse{Products: links}, nil
}

func (s BackofficeService) UpdateCompetitorPrice(ctx context.Context, req productparam.UpdateCompetitorPriceRequest) (productparam.UpdateCompetitorPriceResponse, error) {
	const op = "productservice.UpdateCompetitorPrice"

	if err := s.repo.SetCompetitorPrice(ctx, req.ProductID, req.MinPrice, req.MaxPrice, req.OffersCount, req.CheckedAt.UTC()); err != nil {
		return productparam.UpdateCompetitorPriceResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.UpdateCompetitorPriceResponse{}, nil
}

// ourPriceAndStock is what a storefront visitor would see: the cheapest
// priced variant, and whether any variant can be sold right now. A product
// that isn't sellable on the website counts as out of stock.
func ourPriceAndStock(p entity.Product) (uint64, bool) {
	var minPrice uint64
	inStock := false
	for _, v := range p.Variants {
		if v.Price > 0 && (minPrice == 0 || v.Price < minPrice) {
			minPrice = v.Price
		}
		if v.CanSell(1) {
			inStock = true
		}
	}
	return minPrice, inStock && p.IsSellableOnWebsite()
}

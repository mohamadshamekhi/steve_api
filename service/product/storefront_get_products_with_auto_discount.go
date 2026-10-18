package productservice

import (
	"context"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetProductsWithAutoDiscount(ctx context.Context, _ productparam.GetProductsWithAutoDiscountsRequest) (productparam.GetProductsWithAutoDiscountsResponse, error) {
	const op = "productservice.GetProductsWithAutoDiscount"
	res, err := s.discountProvider.GetProductIDsThatHaveActiveAutoDiscount(ctx, discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest{})
	if err != nil {
		return productparam.GetProductsWithAutoDiscountsResponse{}, richerror.New(op).WithWarpError(err)
	}
	var productList []productparam.StorefrontProductList
	for pid, d := range res.DiscountMap {
		prod, err := s.repo.FindByID(ctx, pid)
		if err != nil {
			return productparam.GetProductsWithAutoDiscountsResponse{}, richerror.New(op).WithWarpError(err)
		}
		productList = append(productList, productparam.StorefrontProductList{
			ID:           prod.ID,
			Name:         prod.Name,
			ENTitle:      prod.ENTitle,
			Slug:         prod.Slug,
			Gallery:      prod.Gallery,
			IsPublished:  prod.IsPublished,
			IsSecondHand: prod.IsSecondHand,
			Discount:     d,
			Variant:      productparam.ToStorefrontVariants(prod.Variants),
			Rating:       prod.Rating,
		})

	}

	return productparam.GetProductsWithAutoDiscountsResponse{
		Products: productList,
	}, err

}

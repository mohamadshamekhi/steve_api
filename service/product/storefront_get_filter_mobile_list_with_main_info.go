package productservice

import (
	"context"

	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetFilterMobileListWithMainInfo(ctx context.Context, req productparam.GetFilterMobileListWithMainInfoRequest) (productparam.GetFilterMobileListWithMainInfoResponse, error) {
	const op = "productservice.GetFilterMobileListWithMainInfo"
	mobiles, err := s.repo.GetFilterMobileListWithMainInfo(ctx, req.Name, req.BrandID)
	if err != nil {
		return productparam.GetFilterMobileListWithMainInfoResponse{}, richerror.New(op).WithWarpError(err)
	}

	var res []productparam.MobileListWithMainInfo
	for _, m := range mobiles {
		res = append(res, productparam.MobileListWithMainInfo{
			ID:           m.ID,
			Name:         m.Name,
			Slug:         m.Slug,
			BrandID:      m.BrandID,
			CategoryID:   m.CategoryID,
			Gallery:      m.Gallery,
			IsPublished:  m.IsPublished,
			IsSecondHand: m.IsSecondHand,
		})
	}

	return productparam.GetFilterMobileListWithMainInfoResponse{
		Products: res,
	}, nil
}

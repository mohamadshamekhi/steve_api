package brandservice

import (
	"context"

	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	BrandReader BrandReader
	mediaSvc    MediaService
}

func NewStorefrontService(br BrandReader, mediaSvc MediaService) StorefrontService {
	return StorefrontService{
		BrandReader: br,
		mediaSvc:    mediaSvc,
	}
}

func (s StorefrontService) GetAllBrands(ctx context.Context, req brandparam.GetAllStorefrontBrandsParamRequest) (brandparam.GetAllStorefrontBrandsParamResponse, error) {
	const op = "brandservice.GetAllBrands"
	resp, err := s.BrandReader.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return brandparam.GetAllStorefrontBrandsParamResponse{}, richerror.New(op).WithWarpError(err)
	}
	var brands []brandparam.BrandInfo

	for _, b := range resp.Brands {
		r, e := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: b.MediaID})
		if e != nil {
			logger.L().Error("warning: could not get media url", "error", e)
		}
		brands = append(brands, brandparam.BrandInfo{
			ID:          b.ID,
			Name:        b.Name,
			Slug:        b.Slug,
			Description: b.Description,
			URL:         r.URL,
			IsFavorite:  b.IsFavorite,
		})
	}

	return brandparam.GetAllStorefrontBrandsParamResponse{Brands: brands, Pagination: struct {
		PageSize   uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
		Total      uint `json:"total"`
	}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)}}, nil
}

package categoryservice

import (
	"context"

	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	repo     CategoryReader
	mediaSvc MediaService
}

func NewStorefrontService(repo CategoryReader, mediaSvc MediaService) StorefrontService {
	return StorefrontService{
		repo:     repo,
		mediaSvc: mediaSvc,
	}
}

func (s StorefrontService) GetAllCategories(ctx context.Context, req categoryparam.GetAllStorefrontCategoriesRequest) (categoryparam.GetAllStorefrontCategoriesResponse, error) {
	const op = "Category.GetAllCategories"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return categoryparam.GetAllStorefrontCategoriesResponse{}, richerror.New(op).WithWarpError(err)
	}
	var categoriesInfo []categoryparam.CategoryInfo

	mediaIDs := make([]string, 0, len(resp.Categories))
	mediaIDSet := make(map[string]struct{})
	for _, c := range resp.Categories {
		if c.MediaID != "" {
			if _, exists := mediaIDSet[c.MediaID]; !exists {
				mediaIDSet[c.MediaID] = struct{}{}
				mediaIDs = append(mediaIDs, c.MediaID)
			}
		}
	}

	urlMap, urlErr := s.mediaSvc.GetURLsByMediaIDs(ctx, mediaparam.GetUrlsByIDsRequest{IDs: mediaIDs})
	if urlErr != nil {
		logger.L().Error("warning: could not get media urls", "error", urlErr)
	}
	for _, c := range resp.Categories {
		categoriesInfo = append(categoriesInfo, categoryparam.CategoryInfo{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			URL:         urlMap.URLs[c.MediaID],
			Level:       c.Level,
			Ancestors:   c.Ancestors,
			Slug:        c.Slug,
		})
	}

	return categoryparam.GetAllStorefrontCategoriesResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Categories: categoriesInfo,
	}, nil
}

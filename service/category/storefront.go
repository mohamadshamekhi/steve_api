package categoryservice

import (
	"context"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	repo            CategoryReader
	mediaSvc        MediaService
	AttrProvider    AttrProvider
	productProvider ProductProvider
}

type ProductProvider interface {
	GetFilterableAttributes(ctx context.Context, req productparam.GetFilterForCategoryRequest) (productparam.GetFilterForCategoryResponse, error)
	GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx context.Context, req productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryRequest) (productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryResponse, error)
}

func NewStorefrontService(repo CategoryReader, mediaSvc MediaService, AttrProvider AttrProvider, productProvider ProductProvider) StorefrontService {
	return StorefrontService{
		repo:            repo,
		mediaSvc:        mediaSvc,
		AttrProvider:    AttrProvider,
		productProvider: productProvider,
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

func (s StorefrontService) GetFilterForCategory(ctx context.Context, req categoryparam.GetFilterForCategoryRequest) (categoryparam.GetFilterForCategoryResponse, error) {
	const op = "Category.GetFilterForCategory"

	category, err := s.repo.FindByID(ctx, req.CategoryID)

	aggregatedOptions, err := s.productProvider.GetFilterableAttributes(ctx, productparam.GetFilterForCategoryRequest{CategoryID: req.CategoryID})
	if err != nil {
		return categoryparam.GetFilterForCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	fmt.Println("aggregatedOptions", aggregatedOptions.AggregatedFilterOption)
	definitionsMap := make(map[string]entity.AttributeDefinition)

	attrRes, err := s.AttrProvider.FindByID(ctx, attributeparam.GetAttributeRequest{ID: category.AttributeSetID})
	if err != nil {
		return categoryparam.GetFilterForCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}

	if category.AttributeSetID != "" {
		for _, def := range attrRes.Attr.Attributes {
			definitionsMap[def.Key] = def
		}
	}

	finalFilters := make([]categoryparam.FilterOption, 0, len(aggregatedOptions.AggregatedFilterOption))
	for _, option := range aggregatedOptions.AggregatedFilterOption {
		if def, ok := definitionsMap[option.Key]; ok {
			finalFilters = append(finalFilters, categoryparam.FilterOption{
				Key:     option.Key,
				Name:    def.Name,
				Unit:    def.Unit,
				Options: option.Values,
			})
		}
	}
	minMaxRes, err := s.productProvider.GetMinAdnMaxPriceOfProductWithSpecialCategory(ctx, productparam.GetMinAdnMaxPriceOfProductWithSpecialCategoryRequest{CategoryID: category.ID})
	if err != nil {
		return categoryparam.GetFilterForCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}

	return categoryparam.GetFilterForCategoryResponse{
		FilterOption: finalFilters,
		MaxPrice:     minMaxRes.MaxPrice,
		MinPrice:     minMaxRes.MinPrice,
	}, nil
}

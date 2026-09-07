package categoryservice

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/abolfazlnorzad/steve/entity"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo         CategoryReadWriter
	mediaSvc     MediaService
	AttrProvider AttrProvider
}

func NewBackofficeService(repo CategoryReadWriter, mediaSvc MediaService, AttrProvider AttrProvider) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc, AttrProvider: AttrProvider}
}

type AttrProvider interface {
	FindByID(ctx context.Context, req attributeparam.GetAttributeRequest) (attributeparam.GetAttributeResponse, error)
}

func (s BackofficeService) Create(ctx context.Context, req categoryparam.CreateCategoryRequest) (categoryparam.CreateCategoryResponse, error) {
	const op = "categoryservice.Create"

	var pid string

	if req.ParentID != nil {
		pid = *req.ParentID
	} else {
		pid = ""
	}

	c := entity.Category{
		Name:           req.Name,
		Description:    req.Description,
		MediaID:        req.MediaID,
		ParentID:       pid,
		AttributeSetID: req.AttributeSetID,
	}

	if req.ParentID != nil {
		p, fErr := s.repo.FindByID(ctx, *req.ParentID)
		if fErr != nil {
			var rErr *richerror.RichError
			if errors.As(fErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return categoryparam.CreateCategoryResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("global.not-found").
					WithMeta(map[string]any{"parent_id": req.ParentID})
			}
			return categoryparam.CreateCategoryResponse{}, richerror.New(op).WithWarpError(fErr)
		}

		c.Ancestors = p.Ancestors
		c.Ancestors = append(c.Ancestors, entity.Ancestor{
			ID:    p.ID,
			Name:  p.Name,
			Slug:  p.Slug,
			Level: p.Level,
		})
		c.Level = p.Level + 1
	} else {
		c.Level = 0
		c.Ancestors = []entity.Ancestor{}
	}

	slug := generateSlug(req.Name)
	counter := 1

	for {
		existingCategory, err := s.repo.FindBySlug(ctx, slug)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				break
			}
			return categoryparam.CreateCategoryResponse{}, richerror.New(op).WithWarpError(err)
		}

		if existingCategory.ID == "" {
			break
		}

		slug = fmt.Sprintf("%s-%d", slug, counter)
		counter++
	}

	c.Slug = slug

	c.BeforeCreate()

	if err := s.repo.Create(ctx, c); err != nil {
		return categoryparam.CreateCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return categoryparam.CreateCategoryResponse{}, nil
}

func (s BackofficeService) Delete(ctx context.Context, req categoryparam.DeleteCategoryRequest) (categoryparam.DeleteCategoryResponse, error) {
	const op = "categoryservice.Delete"

	category, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return categoryparam.DeleteCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}

	if !category.DeletedAt.IsZero() {
		return categoryparam.DeleteCategoryResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	if dErr := s.repo.Delete(ctx, category); dErr != nil {
		return categoryparam.DeleteCategoryResponse{}, richerror.New(op).WithWarpError(dErr)
	}
	// todo-write background job for delete media after one day and children's media
	//go func() {
	//	if _, err := s.mediaSvc.DeleteMedia(context.Background(), mediaparam.DeleteMediaRequest{ID: brand.MediaID}); err != nil {
	//		logger.L().Error("ERROR: background media deletion failed for ID %s: %v\n", brand.MediaID, err)
	//	}
	//}()

	return categoryparam.DeleteCategoryResponse{}, nil
}

func (s BackofficeService) GetCategory(ctx context.Context, req categoryparam.GetCategoryRequest) (categoryparam.GetCategoryResponse, error) {
	const op = "categoryservice.GetCategory"
	category, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return categoryparam.GetCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !category.DeletedAt.IsZero() {
		return categoryparam.GetCategoryResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	m, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: category.MediaID})
	if err != nil {
		logger.L().Error(fmt.Sprintf("failed to get media url by id: %s - err:%v", category.MediaID, err))
	}
	r, err := s.AttrProvider.FindByID(ctx, attributeparam.GetAttributeRequest{ID: category.AttributeSetID})
	if err != nil {
		return categoryparam.GetCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return categoryparam.GetCategoryResponse{
		Category: categoryparam.CategoryInfo{
			ID:             category.ID,
			Name:           category.Name,
			Description:    category.Description,
			URL:            m.URL,
			ParentID:       strPtrOrNil(category.ParentID),
			Ancestors:      category.Ancestors,
			Level:          category.Level,
			Slug:           category.Slug,
			AttributeSetID: strPtrOrNil(category.AttributeSetID),
			AttributeSet:   r.Attr.Attributes,
		},
	}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req categoryparam.GetAllCategoriesRequest) (categoryparam.GetAllCategoriesResponse, error) {
	const op = "categoryservice.GetAll"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return categoryparam.GetAllCategoriesResponse{}, richerror.New(op).WithWarpError(err)
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
		r, _ := s.AttrProvider.FindByID(ctx, attributeparam.GetAttributeRequest{ID: c.AttributeSetID})

		categoriesInfo = append(categoriesInfo, categoryparam.CategoryInfo{
			ID:             c.ID,
			Name:           c.Name,
			Description:    c.Description,
			URL:            urlMap.URLs[c.MediaID],
			Level:          c.Level,
			ParentID:       strPtrOrNil(c.ParentID),
			Ancestors:      c.Ancestors,
			Slug:           c.Slug,
			AttributeSetID: strPtrOrNil(c.AttributeSetID),
			AttributeSet:   r.Attr.Attributes,
		})
	}

	return categoryparam.GetAllCategoriesResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Categories: categoriesInfo,
	}, nil
}

func (s BackofficeService) Update(ctx context.Context, req categoryparam.UpdateCategoryRequest) (categoryparam.UpdateCategoryResponse, error) {
	const op = "categoryservice.Update"
	category, err := s.repo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return categoryparam.UpdateCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !category.DeletedAt.IsZero() {
		return categoryparam.UpdateCategoryResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}
	var oldMediaIDToDelete string

	if req.Name != nil && *req.Name != category.Name {
		category.Name = *req.Name

		slug := generateSlug(*req.Name)
		counter := 1

		for {
			existingCategory, err := s.repo.FindBySlug(ctx, slug)
			if err != nil {
				var rErr *richerror.RichError
				if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
					break
				}
				return categoryparam.UpdateCategoryResponse{}, richerror.New(op).WithWarpError(err)
			}

			if existingCategory.ID == "" {
				break
			}

			slug = fmt.Sprintf("%s-%d", slug, counter)
			counter++
		}

		category.Slug = slug
	}

	if req.Description != nil {
		category.Description = *req.Description
	}

	if req.MediaID != nil && *req.MediaID != category.MediaID {

		_, mediaErr := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: *req.MediaID})
		if mediaErr != nil {
			var rErr *richerror.RichError
			if errors.As(mediaErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return categoryparam.UpdateCategoryResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("media.not_found").
					WithMeta(map[string]any{"media_id": *req.MediaID})
			}
			return categoryparam.UpdateCategoryResponse{}, richerror.New(op).WithWarpError(mediaErr)
		}

		oldMediaIDToDelete = category.MediaID

		category.MediaID = *req.MediaID
	}

	if req.AttributeSetID != nil {
		category.AttributeSetID = *req.AttributeSetID
	}

	if req.ParentID != nil && *req.ParentID != category.ParentID {
		newParentID := *req.ParentID

		if newParentID == "" {
			category.ParentID = ""
			category.Level = 0
			category.Ancestors = []entity.Ancestor{}
		} else {
			if newParentID == category.ID {
				return categoryparam.UpdateCategoryResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("global.validation_failed").
					WithMeta(map[string]any{"parent_id": newParentID})
			}

			p, fErr := s.repo.FindByID(ctx, newParentID)
			if fErr != nil {
				var rErr *richerror.RichError
				if errors.As(fErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
					return categoryparam.UpdateCategoryResponse{}, richerror.New(op).
						WithKind(richerror.KindInvalid).
						WithMessageKey("global.not-found").
						WithMeta(map[string]any{"parent_id": newParentID})
				}
				return categoryparam.UpdateCategoryResponse{}, richerror.New(op).WithWarpError(fErr)
			}

			// guard against making a descendant the new parent (would create a cycle)
			for _, a := range p.Ancestors {
				if a.ID == category.ID {
					return categoryparam.UpdateCategoryResponse{}, richerror.New(op).
						WithKind(richerror.KindInvalid).
						WithMessageKey("global.validation_failed").
						WithMeta(map[string]any{"parent_id": newParentID})
				}
			}

			ancestors := make([]entity.Ancestor, 0, len(p.Ancestors)+1)
			ancestors = append(ancestors, p.Ancestors...)
			ancestors = append(ancestors, entity.Ancestor{
				ID:    p.ID,
				Name:  p.Name,
				Slug:  p.Slug,
				Level: p.Level,
			})

			category.ParentID = newParentID
			category.Ancestors = ancestors
			category.Level = p.Level + 1
		}
	}

	if uErr := s.repo.Update(ctx, category); uErr != nil {
		return categoryparam.UpdateCategoryResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
	}
	if oldMediaIDToDelete != "" {
		go func() {
			r, err := s.mediaSvc.FindMediaByID(context.Background(), mediaparam.GetMediaByIDRequest{ID: oldMediaIDToDelete})
			if err != nil {
				logger.L().Error("error in find old media", oldMediaIDToDelete)
			}
			if _, err := s.mediaSvc.DeleteMedia(context.Background(), mediaparam.DeleteMediaRequest{Media: r.Media}); err != nil {
				logger.L().Error("ERROR: background media deletion failed for ID %s: %v\n", oldMediaIDToDelete, err)
			}
		}()
	}

	return categoryparam.UpdateCategoryResponse{}, nil
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func generateSlug(name string) string {
	slug := strings.ReplaceAll(name, " ", "-")
	slug = strings.ReplaceAll(slug, "‌", "-")

	re := regexp.MustCompile(`[^\p{L}\p{N}-]+`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile("-{2,}")
	slug = re.ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

func (s BackofficeService) GetCategoriesByIDs(ctx context.Context, req categoryparam.GetCategoriesByIDsRequest) (categoryparam.GetCategoriesByIDsResponse, error) {
	const op = "categoryservice.GetCategoriesByIDs"

	categories, err := s.repo.GetCategoriesByIDs(ctx, req.IDs)
	if err != nil {
		return categoryparam.GetCategoriesByIDsResponse{}, richerror.New(op).WithWarpError(err)
	}
	return categoryparam.GetCategoriesByIDsResponse{Categories: categories}, nil
}

func (s BackofficeService) GetCategoryBySlug(ctx context.Context, req categoryparam.GetCategoryBySlugRequest) (categoryparam.GetCategoryBySlugResponse, error) {
	const op = "categoryservice.GetCategoryBySlug"

	resp, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil {
		return categoryparam.GetCategoryBySlugResponse{}, richerror.New(op).WithWarpError(err)
	}

	category := entity.Category{
		ID:             resp.ID,
		Name:           resp.Name,
		Slug:           resp.Slug,
		Description:    resp.Description,
		MediaID:        resp.MediaID,
		ParentID:       resp.ParentID,
		Ancestors:      resp.Ancestors,
		Level:          resp.Level,
		AttributeSetID: resp.AttributeSetID,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
		DeletedAt:      resp.DeletedAt,
	}

	return categoryparam.GetCategoryBySlugResponse{Category: category}, nil
}

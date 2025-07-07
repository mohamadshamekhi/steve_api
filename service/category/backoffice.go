package categoryservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/abolfazlnorzad/steve/entity"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"regexp"
	"strings"
)

type BackofficeService struct {
	repo     CategoryReadWriter
	mediaSvc MediaService
}

func NewBackofficeService(repo CategoryReadWriter, mediaSvc MediaService) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc}
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
		Name:        req.Name,
		Description: req.Description,
		MediaID:     req.MediaID,
		ParentID:    pid,
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

	m, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: category.MediaID})
	if err != nil {
		return categoryparam.GetCategoryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return categoryparam.GetCategoryResponse{
		Category: categoryparam.CategoryInfo{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			URL:         m.URL,
			Ancestors:   category.Ancestors,
			Level:       category.Level,
			Slug:        category.Slug,
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

	var mIDS []string
	for _, c := range resp.Categories {
		mIDS = append(mIDS, c.MediaID)
	}

	r, err := s.mediaSvc.GetURLsByMediaIDs(ctx, mediaparam.GetUrlsByIDsRequest{IDs: mIDS})
	if err != nil {
		return categoryparam.GetAllCategoriesResponse{}, err
	}
	for _, c := range resp.Categories {
		categoriesInfo = append(categoriesInfo, categoryparam.CategoryInfo{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			URL:         r.URLs[c.ID],
			Level:       c.Level,
			Ancestors:   c.Ancestors,
			Slug:        c.Slug,
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

func generateSlug(name string) string {
	slug := strings.ReplaceAll(name, " ", "-")
	slug = strings.ReplaceAll(slug, "‌", "-")

	re := regexp.MustCompile(`[^\p{L}\p{N}-]+`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile("-{2,}")
	slug = re.ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

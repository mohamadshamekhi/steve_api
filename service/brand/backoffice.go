package brandservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/abolfazlnorzad/steve/pkg/slug"
)

type BackofficeService struct {
	repo     BrandReadWriter
	mediaSvc MediaService
}

func NewBackofficeService(repo BrandReadWriter, mediaSvc MediaService) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc}
}

func (s BackofficeService) Create(ctx context.Context, req brandparam.CreateBrandRequest) (brandparam.CreateBrandResponse, error) {
	const op = "brandservice.Create"

	brandSlug, sErr := s.uniqueSlug(ctx, req.Slug, req.Name, "")
	if sErr != nil {
		return brandparam.CreateBrandResponse{}, richerror.New(op).WithWarpError(sErr)
	}

	brand := entity.Brand{
		Name:        req.Name,
		Slug:        brandSlug,
		Description: req.Description,
		MediaID:     req.MediaID,
		IsFavorite:  req.IsFavorite,
	}
	brand.BeforeCreate()

	if err := s.repo.Create(ctx, brand); err != nil {
		return brandparam.CreateBrandResponse{}, richerror.New(op).WithWarpError(err)
	}

	return brandparam.CreateBrandResponse{}, nil
}

// uniqueSlug builds a unique slug for a brand. If explicit is non-empty it is
// normalized and used as the base (letting an admin pin an English slug like
// "samsung" for a Persian brand name); otherwise the base is derived from
// name, mirroring category's auto-slug behavior. A numeric suffix is
// appended until the slug doesn't collide with another brand; excludeID lets
// an update keep its own current slug instead of "colliding with itself".
func (s BackofficeService) uniqueSlug(ctx context.Context, explicit, name, excludeID string) (string, error) {
	base := slug.GenerateSlug(explicit)
	if base == "" {
		base = slug.GenerateSlug(name)
	}

	candidate := base
	counter := 1
	for {
		existing, err := s.repo.FindBySlug(ctx, candidate)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return candidate, nil
			}
			return "", err
		}
		if existing.ID == "" || existing.ID == excludeID {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}
}

func (s BackofficeService) GetBrand(ctx context.Context, req brandparam.GetBrandRequest) (brandparam.GetBrandResponse, error) {
	const op = "brandservice.GetBrand"

	brand, err := s.repo.FindByID(ctx, req.ID)

	if err != nil {
		return brandparam.GetBrandResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !brand.DeletedAt.IsZero() {
		return brandparam.GetBrandResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}
	resp, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: brand.MediaID})
	if err != nil {
		logger.L().Error("failed to get url brand , brandID :", brand.ID, err)
	}
	return brandparam.GetBrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
		Url:         resp.URL,
		IsFavorite:  brand.IsFavorite,
	}, nil
}

func (s BackofficeService) Update(ctx context.Context, req brandparam.UpdateBrandRequest) (brandparam.UpdateBrandResponse, error) {
	const op = "brandservice.Update"
	brand, fErr := s.repo.FindByID(ctx, req.BrandID)
	if fErr != nil {
		return brandparam.UpdateBrandResponse{}, richerror.New(op).WithWarpError(fErr)
	}

	if !brand.DeletedAt.IsZero() {
		return brandparam.UpdateBrandResponse{}, richerror.New(op).
			WithKind(richerror.KindNotFound).
			WithMessageKey("global.not-found")
	}

	var oldMediaIDToDelete string

	nameChanged := req.Name != nil && *req.Name != brand.Name
	slugRequested := req.Slug != nil && *req.Slug != ""

	if req.Name != nil {
		brand.Name = *req.Name
	}

	if slugRequested || nameChanged {
		var explicit string
		if slugRequested {
			explicit = *req.Slug
		}
		newSlug, sErr := s.uniqueSlug(ctx, explicit, brand.Name, brand.ID)
		if sErr != nil {
			return brandparam.UpdateBrandResponse{}, richerror.New(op).WithWarpError(sErr)
		}
		brand.Slug = newSlug
	}

	if req.Description != nil {
		brand.Description = *req.Description
	}

	if req.MediaID != nil && *req.MediaID != brand.MediaID {

		_, mediaErr := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: *req.MediaID})
		if mediaErr != nil {
			var rErr *richerror.RichError
			if errors.As(mediaErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return brandparam.UpdateBrandResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("media.not_found").
					WithMeta(map[string]any{"media_id": *req.MediaID})
			}
			return brandparam.UpdateBrandResponse{}, richerror.New(op).WithWarpError(mediaErr)
		}

		oldMediaIDToDelete = brand.MediaID

		brand.MediaID = *req.MediaID
	}

	if req.IsFavorite != nil {
		brand.IsFavorite = *req.IsFavorite
	}

	if uErr := s.repo.Update(ctx, brand); uErr != nil {
		return brandparam.UpdateBrandResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
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

	return brandparam.UpdateBrandResponse{}, nil

}

func (s BackofficeService) Delete(ctx context.Context, req brandparam.DeleteBrandRequest) (brandparam.DeleteBrandResponse, error) {
	const op = "brandservice.Delete"

	brand, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return brandparam.DeleteBrandResponse{}, richerror.New(op).WithWarpError(err)
	}

	if !brand.DeletedAt.IsZero() {
		return brandparam.DeleteBrandResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	if dErr := s.repo.Delete(ctx, brand); dErr != nil {
		return brandparam.DeleteBrandResponse{}, richerror.New(op).WithWarpError(dErr)
	}
	// todo-write background job for delete media after one day
	//go func() {
	//	if _, err := s.mediaSvc.DeleteMedia(context.Background(), mediaparam.DeleteMediaRequest{ID: brand.MediaID}); err != nil {
	//		logger.L().Error("ERROR: background media deletion failed for ID %s: %v\n", brand.MediaID, err)
	//	}
	//}()

	return brandparam.DeleteBrandResponse{}, nil
}

func (s BackofficeService) All(ctx context.Context, req brandparam.GetAllBrandsParamRequest) (brandparam.GetAllBrandsParamResponse, error) {
	const op = "brandservice.All"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return brandparam.GetAllBrandsParamResponse{}, richerror.New(op).WithWarpError(err)
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

	return brandparam.GetAllBrandsParamResponse{Brands: brands, Pagination: struct {
		PageSize   uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
		Total      uint `json:"total"`
	}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)}}, nil

}

func (s BackofficeService) GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error) {
	const op = "brandservice.GetBrandsByIDs"

	res, err := s.repo.GetBrandsByIDs(ctx, req.IDs)
	if err != nil {
		return brandparam.GetBrandsByIDsResponse{}, err
	}

	return brandparam.GetBrandsByIDsResponse{
		Brands: res,
	}, nil

}

package brandservice

import (
	"context"
	"errors"
	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
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

	brand := entity.Brand{
		Name:        req.Name,
		Description: req.Description,
		MediaID:     req.MediaID,
	}
	brand.BeforeCreate()

	if err := s.repo.Create(ctx, brand); err != nil {
		return brandparam.CreateBrandResponse{}, richerror.New(op).WithWarpError(err)
	}

	return brandparam.CreateBrandResponse{}, nil
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
		return brandparam.GetBrandResponse{}, richerror.New(op).WithWarpError(err)
	}
	return brandparam.GetBrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Description: brand.Description,
		Url:         resp.URL,
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

	if req.Name != nil {
		brand.Name = *req.Name
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

package settingservice

import (
	"context"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BannerBackofficeService struct {
	repo     BannerRepository
	mediaSvc MediaService
}

type BannerRepository interface {
	Create(ctx context.Context, b entity.Banner) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (entity.Banner, error)
	Update(ctx context.Context, b *entity.Banner) error
	List(ctx context.Context, isActive *bool) ([]entity.Banner, error)
}
type MediaService interface {
	GetURLByMediaID(ctx context.Context, req mediaparam.GetUrlByIDRequest) (mediaparam.GetUrlByIDResponse, error)
	GetURLsByMediaIDs(ctx context.Context, req mediaparam.GetUrlsByIDsRequest) (mediaparam.GetUrlsByIDsResponse, error)
	DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error)
	FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error)
}

func NewBannerBackofficeService(repo BannerRepository, mediaSvc MediaService) BannerBackofficeService {
	return BannerBackofficeService{repo: repo, mediaSvc: mediaSvc}
}

func (bs BannerBackofficeService) Create(ctx context.Context, req settingparam.CreateBannerRequest) (settingparam.CreateBannerResponse, error) {
	const op = "settingservice.BannerService.Create"
	b := entity.Banner{
		Title:     req.Title,
		MediaID:   req.MediaID,
		Link:      req.Link,
		Location:  req.Location,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		AltText:   req.AltText,
		StartsAt:  req.StartsAt,
		ExpiresAt: req.ExpiresAt,
	}
	err := bs.repo.Create(ctx, b)
	if err != nil {
		return settingparam.CreateBannerResponse{}, richerror.New(op).WithWarpError(err)
	}

	return settingparam.CreateBannerResponse{}, nil
}

func (bs BannerBackofficeService) Delete(ctx context.Context, req settingparam.DeleteBannerRequest) (brandparam.DeleteBrandResponse, error) {
	const op = "settingservice.BannerService.Delete"
	err := bs.repo.Delete(ctx, req.ID)
	if err != nil {
		return brandparam.DeleteBrandResponse{}, richerror.New(op).WithWarpError(err)
	}
	return brandparam.DeleteBrandResponse{}, nil
}

func (bs BannerBackofficeService) Update(ctx context.Context, req settingparam.UpdateBannerRequest) (settingparam.UpdateBannerResponse, error) {
	const op = "settingservice.BannerService.Update"
	b, err := bs.repo.FindByID(ctx, req.ID)
	if err != nil {
		return settingparam.UpdateBannerResponse{}, richerror.New(op).WithWarpError(err)
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	if req.MediaID != nil {
		b.MediaID = *req.MediaID
	}
	if req.Location != nil {
		b.Location = *req.Location
	}
	if req.SortOrder != nil {
		b.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		b.IsActive = *req.IsActive
	}
	if req.AltText != nil {
		b.AltText = *req.AltText
	}
	if req.StartsAt != nil {
		b.StartsAt = req.StartsAt
	}
	if req.ExpiresAt != nil {
		b.ExpiresAt = req.ExpiresAt
	}
	if req.Link != nil {
		b.Link = *req.Link
	}
	err = bs.repo.Update(ctx, &b)
	if err != nil {
		return settingparam.UpdateBannerResponse{}, richerror.New(op).WithWarpError(err)
	}
	return settingparam.UpdateBannerResponse{}, nil
}

func (bs BannerBackofficeService) List(ctx context.Context, req settingparam.ListBannerRequest) (settingparam.ListBannerResponse, error) {
	const op = "settingservice.BannerService.List"
	banners, err := bs.repo.List(ctx, req.IsActive)
	if err != nil {
		return settingparam.ListBannerResponse{}, richerror.New(op).WithWarpError(err)
	}
	var result []settingparam.BannerInfo
	var mIds []string
	for _, banner := range banners {
		mIds = append(mIds, banner.MediaID)
	}

	resp, err := bs.mediaSvc.GetURLsByMediaIDs(ctx, mediaparam.GetUrlsByIDsRequest{IDs: mIds})
	if err != nil {
		fmt.Println(richerror.New(op).WithWarpError(err).Error())
	}

	for _, b := range banners {

		result = append(result, settingparam.BannerInfo{
			ID:        b.ID,
			Title:     b.Title,
			MediaID:   b.MediaID,
			ImageURL:  resp.URLs[b.MediaID],
			Link:      b.Link,
			Location:  b.Location,
			SortOrder: b.SortOrder,
			IsActive:  b.IsActive,
			AltText:   b.AltText,
			CreatedAt: b.CreatedAt,
		})
	}
	return settingparam.ListBannerResponse{
		Banners: result,
	}, nil
}

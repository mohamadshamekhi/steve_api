package pricereviewservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	pricereviewparam "github.com/abolfazlnorzad/steve/param/pricereview"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type SettingRepository interface {
	Get(ctx context.Context) (entity.PriceReviewSetting, error)
	Update(ctx context.Context, intervalDays int) error
}

type ProductRepository interface {
	FindNeedingPriceReview(ctx context.Context, cutoff time.Time, offset, limit uint) ([]entity.Product, int64, error)
	MarkPriceReviewed(ctx context.Context, productID string, at time.Time) error
}

type BackofficeService struct {
	settingRepo SettingRepository
	productRepo ProductRepository
}

func NewBackofficeService(settingRepo SettingRepository, productRepo ProductRepository) BackofficeService {
	return BackofficeService{settingRepo: settingRepo, productRepo: productRepo}
}

func (s BackofficeService) GetSetting(ctx context.Context, _ pricereviewparam.GetSettingRequest) (pricereviewparam.GetSettingResponse, error) {
	const op = "pricereviewservice.GetSetting"

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return pricereviewparam.GetSettingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return pricereviewparam.GetSettingResponse{IntervalDays: setting.IntervalDays}, nil
}

func (s BackofficeService) UpdateSetting(ctx context.Context, req pricereviewparam.UpdateSettingRequest) (pricereviewparam.UpdateSettingResponse, error) {
	const op = "pricereviewservice.UpdateSetting"

	if req.IntervalDays <= 0 {
		return pricereviewparam.UpdateSettingResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input")
	}
	if err := s.settingRepo.Update(ctx, req.IntervalDays); err != nil {
		return pricereviewparam.UpdateSettingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return pricereviewparam.UpdateSettingResponse{}, nil
}

func (s BackofficeService) ListNeedsReview(ctx context.Context, req pricereviewparam.ListNeedsReviewRequest) (pricereviewparam.ListNeedsReviewResponse, error) {
	const op = "pricereviewservice.ListNeedsReview"

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return pricereviewparam.ListNeedsReviewResponse{}, richerror.New(op).WithWarpError(err)
	}

	cutoff := time.Now().AddDate(0, 0, -setting.IntervalDays)
	products, total, err := s.productRepo.FindNeedingPriceReview(ctx, cutoff, req.Pagination.GetOffset(), req.Pagination.GetPageSize())
	if err != nil {
		return pricereviewparam.ListNeedsReviewResponse{}, richerror.New(op).WithWarpError(err)
	}

	items := make([]pricereviewparam.NeedsReviewItem, 0, len(products))
	for _, p := range products {
		lastReviewed := p.PriceReviewedAt
		if lastReviewed.IsZero() {
			lastReviewed = p.CreatedAt
		}
		daysOverdue := int(time.Since(lastReviewed).Hours()/24) - setting.IntervalDays
		if daysOverdue < 0 {
			daysOverdue = 0
		}
		items = append(items, pricereviewparam.NeedsReviewItem{
			ProductID:      p.ID,
			ProductName:    p.Name,
			CurrentPrice:   p.DisplayPrice(),
			LastReviewedAt: lastReviewed,
			DaysOverdue:    daysOverdue,
		})
	}

	return pricereviewparam.ListNeedsReviewResponse{Items: items, Total: total}, nil
}

func (s BackofficeService) Confirm(ctx context.Context, req pricereviewparam.ConfirmRequest) (pricereviewparam.ConfirmResponse, error) {
	const op = "pricereviewservice.Confirm"

	if req.ProductID == "" {
		return pricereviewparam.ConfirmResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input")
	}
	if err := s.productRepo.MarkPriceReviewed(ctx, req.ProductID, time.Now()); err != nil {
		return pricereviewparam.ConfirmResponse{}, richerror.New(op).WithWarpError(err)
	}
	return pricereviewparam.ConfirmResponse{}, nil
}

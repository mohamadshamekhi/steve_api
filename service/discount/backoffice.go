package discountservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
)

type BackofficeService struct {
	repo DiscountReadWriter
}

func NewBackofficeService() BackofficeService {
	return BackofficeService{}
}

func (s BackofficeService) Create(ctx context.Context, req discountparam.AddDiscountRequest) (discountparam.AddDiscountResponse, error) {
	const op = "discountservice.Create"

	d := entity.Discount{
		Code:                  req.Code,
		Type:                  entity.DiscountType(req.Type),
		Value:                 req.Value,
		MinimumPurchaseAmount: req.MinimumPurchaseAmount,
		StartsAt:              req.StartsAt,
		ExpiresAt:             req.ExpiresAt,
		TotalUsageLimit:       req.TotalUsageLimit,
		ApplicableProductIDs:  req.ApplicableProductIDs,
		ApplicableCategoryIDs: req.ApplicableCategoryIDs,
		ApplicableUserIDs:     req.ApplicableUserIDs,
		IsActive:              req.IsActive,
	}
	d.BeforeCreate()

	if err := s.repo.Create(ctx, d); err != nil {
		return discountparam.AddDiscountResponse{}, err
	}
	return discountparam.AddDiscountResponse{}, nil
}

func (s BackofficeService) Update(ctx context.Context, req discountparam.UpdateDiscountRequest) (discountparam.UpdateDiscountResponse, error) {
	const op = "discountservice.Update"
	discount, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return discountparam.UpdateDiscountResponse{}, err
	}

	if req.Value != nil {
		discount.Value = *req.Value
	}
	if req.MinimumPurchaseAmount != nil {
		discount.MinimumPurchaseAmount = *req.MinimumPurchaseAmount
	}
	if req.ExpiresAt != nil {
		discount.ExpiresAt = *req.ExpiresAt
	}
	if req.TotalUsageLimit != nil {
		discount.TotalUsageLimit = *req.TotalUsageLimit
	}
	if req.UsageLimitPerUser != nil {
		discount.UsageLimitPerUser = *req.UsageLimitPerUser
	}
	if req.ApplicableProductIDs != nil {
		discount.ApplicableProductIDs = *req.ApplicableProductIDs
	}
	if req.ApplicableCategoryIDs != nil {
		discount.ApplicableCategoryIDs = *req.ApplicableCategoryIDs
	}
	if req.ApplicableUserIDs != nil {
		discount.ApplicableUserIDs = *req.ApplicableUserIDs
	}
	if req.IsActive != nil {
		discount.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, discount); err != nil {
		return discountparam.UpdateDiscountResponse{}, err
	}

	return discountparam.UpdateDiscountResponse{}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req discountparam.GetAllDiscountsParamRequest) (discountparam.GetAllDiscountsParamResponse, error) {
	const op = "backofficeservice.GetAll"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return discountparam.GetAllDiscountsParamResponse{}, err
	}

	return discountparam.GetAllDiscountsParamResponse{Discounts: resp.Discounts, Pagination: struct {
		PageSize   uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
		Total      uint `json:"total"`
	}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)}}, nil
}

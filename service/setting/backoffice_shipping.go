package settingservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeShippingService struct {
	repo ShippingReaderWriter
}

func NewBackofficeShippingService(repo ShippingReaderWriter) BackofficeShippingService {
	return BackofficeShippingService{
		repo: repo,
	}
}

func (s BackofficeShippingService) Create(ctx context.Context, req settingparam.AddShippingRequest) (settingparam.AddShippingResponse, error) {
	const op = "BackofficeShippingService.Create"
	shipping := entity.Shipping{
		Name:        req.Name,
		Description: req.Description,
		Cost:        req.Cost,
		IsActive:    req.IsActive,
	}
	shipping.BeforeCreate()
	if err := s.repo.Create(ctx, shipping); err != nil {
		return settingparam.AddShippingResponse{}, err
	}
	return settingparam.AddShippingResponse{}, nil
}

func (s BackofficeShippingService) Update(ctx context.Context, req settingparam.UpdateShippingRequest) (settingparam.UpdateShippingResponse, error) {
	const op = "BackofficeShippingService.Update"
	shipping, err := s.repo.FindByID(ctx, req.ShippingID)
	if err != nil {
		return settingparam.UpdateShippingResponse{}, richerror.New(op).WithWarpError(err)
	}
	if req.Name != nil {
		shipping.Name = *req.Name
	}

	if req.Cost != nil {
		shipping.Cost = *req.Cost
	}
	if req.Description != nil {
		shipping.Description = *req.Description
	}
	if req.IsActive != nil {
		shipping.IsActive = *req.IsActive
	}
	if err = s.repo.Update(ctx, shipping); err != nil {
		return settingparam.UpdateShippingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return settingparam.UpdateShippingResponse{}, nil
}

func (s BackofficeShippingService) FindByID(ctx context.Context, req settingparam.FindByIDShippingRequest) (settingparam.FindByIDShippingResponse, error) {
	const op = "BackofficeShippingService.FindByID"
	shipping, err := s.repo.FindByID(ctx, req.ShippingID)
	if err != nil {
		return settingparam.FindByIDShippingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return settingparam.FindByIDShippingResponse{
		Shipping: shipping,
	}, nil
}

func (s BackofficeShippingService) GetAll(ctx context.Context, req settingparam.GetAllShippingRequest) (settingparam.GetAllShippingResponse, error) {
	const op = "BackofficeShippingService.GetAll"
	shippings, err := s.repo.GetAll(ctx, req.Filter, req.Pagination)
	if err != nil {
		return settingparam.GetAllShippingResponse{}, richerror.New(op).WithWarpError(err)
	}
	return settingparam.GetAllShippingResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(shippings.TotalCount),
		},
		Shipping: shippings.Shipping,
	}, nil
}

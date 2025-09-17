package repairservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	repairparam "github.com/abolfazlnorzad/steve/param/repair"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	Repo      StorefrontRepo
	uProvider UserProvider
	bProvider BrandProvider
	pProvider ProductProvider
}

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error)
}

type BrandProvider interface {
	GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error)
}

type ProductProvider interface {
	GetProductsByIDs(ctx context.Context, req productparam.GetProductsByIDsRequest) (productparam.GetProductsByIDsResponse, error)
}

type StorefrontRepo interface {
	Add(ctx context.Context, repair entity.Repair) error
	All(ctx context.Context, filter commonparam.FilterRequest, pagination commonparam.PaginationRequest) (ListResult, error)
}

func NewStorefrontService(r StorefrontRepo, uProvider UserProvider,
	bProvider BrandProvider,
	pProvider ProductProvider) StorefrontService {
	return StorefrontService{
		Repo:      r,
		uProvider: uProvider,
		bProvider: bProvider,
		pProvider: pProvider,
	}
}

func (s StorefrontService) Create(ctx context.Context, req repairparam.AddRepairRequest) (repairparam.AddRepairResponse, error) {
	const op = "RepairService.Create"
	repair := entity.Repair{
		UserID:      req.UserID,
		BrandID:     req.BrandID,
		ProductID:   req.ProductID,
		Description: req.Description,
		Status:      entity.StatusNew,
		Galleries:   req.Galleries,
	}
	repair.BeforeCreate()
	err := s.Repo.Add(ctx, repair)
	if err != nil {
		return repairparam.AddRepairResponse{}, richerror.New(op).WithWarpError(err)
	}

	return repairparam.AddRepairResponse{}, nil
}

func (s StorefrontService) GetAllRepair(ctx context.Context, req repairparam.GetAllMyRepairParamRequest) (repairparam.GetAllMyRepairParamResponse, error) {
	const op = "StorefrontService.GetAllRepair"
	resp, err := s.Repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return repairparam.GetAllMyRepairParamResponse{}, richerror.New(op).WithWarpError(err)
	}
	var userIDs []string
	for _, repair := range resp.Repairs {
		userIDs = append(userIDs, repair.UserID)
	}
	uResp, err := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
	if err != nil {
		return repairparam.GetAllMyRepairParamResponse{}, err
	}
	var brandIDs []string
	for _, repair := range resp.Repairs {
		brandIDs = append(brandIDs, repair.BrandID)
	}
	bResp, err := s.bProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: brandIDs})
	if err != nil {
		logger.L().Error("can't find brand- err :", err)
	}
	var productIDs []string
	for _, repair := range resp.Repairs {
		productIDs = append(productIDs, repair.ProductID)
	}
	pResp, err := s.pProvider.GetProductsByIDs(ctx, productparam.GetProductsByIDsRequest{IDs: productIDs})
	if err != nil {
		logger.L().Error("can't find product- err :", err)
	}
	var repairInfos []repairparam.RepairInfo
	for _, r := range resp.Repairs {
		var user userparam.UserInfo
		for _, u := range uResp.Users {
			if u.ID == r.UserID {
				user = u
			}
		}
		repairInfos = append(repairInfos, repairparam.RepairInfo{
			ID: r.ID,
			User: struct {
				Name        string `json:"name"`
				PhoneNumber string `json:"phone_number"`
			}{
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			BrandName:   bResp.Brands[r.BrandID].Name,
			ProductName: pResp.Products[r.ProductID].Name,
			Description: r.Description,
			Status:      r.Status,
			Galleries:   r.Galleries,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return repairparam.GetAllMyRepairParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Repairs: repairInfos,
	}, nil
}

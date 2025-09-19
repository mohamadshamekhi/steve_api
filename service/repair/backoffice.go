package repairservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	repairparam "github.com/abolfazlnorzad/steve/param/repair"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo      Repository
	uProvider UserProvider
	bProvider BrandProvider
	pProvider ProductProvider
}

func NewBackofficeService(
	repo Repository,
	uProvider UserProvider,
	bProvider BrandProvider,
	pProvider ProductProvider) BackofficeService {
	return BackofficeService{
		repo:      repo,
		uProvider: uProvider,
		bProvider: bProvider,
		pProvider: pProvider,
	}
}

type Repository interface {
	FindByID(ctx context.Context, id string) (entity.Repair, error)
	All(ctx context.Context, filter commonparam.FilterRequest, pagination commonparam.PaginationRequest) (ListResult, error)
	Update(ctx context.Context, repair entity.Repair) error
}

func (s BackofficeService) GetRepair(ctx context.Context, req repairparam.GetRepairRequest) (repairparam.GetRepairResponse, error) {
	const op = "BackofficeService.GetRepair"
	r, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return repairparam.GetRepairResponse{}, richerror.New(op).WithWarpError(err)
	}
	bResp, bErr := s.bProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: []string{r.BrandID}})
	if bErr != nil {
		logger.L().Error("Failed to get repair brand info", bErr)
	}
	pResp, pErr := s.pProvider.GetProductsByIDs(ctx, productparam.GetProductsByIDsRequest{IDs: []string{r.ProductID}})
	if pErr != nil {
		logger.L().Error("Failed to get repair products", pErr)
	}
	uResp, uErr := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: []string{r.UserID}})
	var u userparam.UserInfo
	for _, user := range uResp.Users {
		u = user
	}
	if uErr != nil {
		logger.L().Error("Failed to get repair users", uErr)
	}
	return repairparam.GetRepairResponse{RepairInfo: repairparam.RepairInfo{
		ID: r.ID,
		User: struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
		},
		BrandName:   bResp.Brands[r.BrandID].Name,
		ProductName: pResp.Products[r.ProductID].Name,
		Description: r.Description,
		Status:      r.Status,
		Galleries:   r.Galleries,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req repairparam.GetAllRepairParamRequest) (repairparam.GetAllRepairParamResponse, error) {
	const op = "BackofficeService.GetAll"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return repairparam.GetAllRepairParamResponse{}, richerror.New(op).WithWarpError(err)
	}
	var userIDs []string
	for _, repair := range resp.Repairs {
		userIDs = append(userIDs, repair.UserID)
	}
	uResp, err := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
	if err != nil {
		return repairparam.GetAllRepairParamResponse{}, err
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

	return repairparam.GetAllRepairParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Repairs: repairInfos,
	}, nil
}

func (s BackofficeService) Update(ctx context.Context, req repairparam.UpdateRepairRequest) (repairparam.UpdateRepairResponse, error) {
	const op = "BackofficeService.Update"
	r, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return repairparam.UpdateRepairResponse{}, richerror.New(op).WithWarpError(err)
	}
	r.Status = req.Status
	r.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, r)
	if err != nil {
		return repairparam.UpdateRepairResponse{}, richerror.New(op).WithWarpError(err)
	}
	return repairparam.UpdateRepairResponse{}, nil
}

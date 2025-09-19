package sellservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
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
	FindByID(ctx context.Context, id string) (entity.Sell, error)
	All(ctx context.Context, filter commonparam.FilterRequest, pagination commonparam.PaginationRequest) (ListResult, error)
	Update(ctx context.Context, repair entity.Sell) error
}

func (s BackofficeService) GetSell(ctx context.Context, req sellparam.GetSellRequest) (sellparam.GetSellResponse, error) {
	const op = "sellservice.GetSell"
	r, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return sellparam.GetSellResponse{}, richerror.New(op).WithWarpError(err)
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
	return sellparam.GetSellResponse{Sell: sellparam.SellInfo{
		ID: r.ID,
		User: struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
		},
		BrandName:         bResp.Brands[r.BrandID].Name,
		ProductName:       pResp.Products[r.ProductID].Name,
		Color:             r.Color,
		BatteryHealth:     r.BatteryHealth,
		RegisterStatus:    r.RegisterStatus,
		PhysicalCondition: r.PhysicalCondition,
		Description:       r.Description,
		Status:            r.Status,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req sellparam.GetAllSellParamRequest) (sellparam.GetAllSellParamResponse, error) {
	const op = "sellservice.GetAll"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return sellparam.GetAllSellParamResponse{}, richerror.New(op).WithWarpError(err)
	}
	var userIDs []string
	for _, sell := range resp.Sells {
		userIDs = append(userIDs, sell.UserID)
	}
	uResp, err := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
	if err != nil {
		return sellparam.GetAllSellParamResponse{}, err
	}
	var brandIDs []string
	for _, repair := range resp.Sells {
		brandIDs = append(brandIDs, repair.BrandID)
	}
	bResp, err := s.bProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: brandIDs})
	if err != nil {
		logger.L().Error("can't find brand- err :", err)
	}
	var productIDs []string
	for _, repair := range resp.Sells {
		productIDs = append(productIDs, repair.ProductID)
	}
	pResp, err := s.pProvider.GetProductsByIDs(ctx, productparam.GetProductsByIDsRequest{IDs: productIDs})
	if err != nil {
		logger.L().Error("can't find product- err :", err)
	}
	var sellInfos []sellparam.SellInfo
	for _, r := range resp.Sells {
		var user userparam.UserInfo
		for _, u := range uResp.Users {
			if u.ID == r.UserID {
				user = u
			}
		}
		sellInfos = append(sellInfos, sellparam.SellInfo{
			ID: r.ID,
			User: struct {
				Name        string `json:"name"`
				PhoneNumber string `json:"phone_number"`
			}{
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			BrandName:         bResp.Brands[r.BrandID].Name,
			ProductName:       pResp.Products[r.ProductID].Name,
			Color:             r.Color,
			BatteryHealth:     r.BatteryHealth,
			RegisterStatus:    r.RegisterStatus,
			PhysicalCondition: r.PhysicalCondition,
			Description:       r.Description,
			Status:            r.Status,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
		})
	}

	return sellparam.GetAllSellParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Sells: sellInfos,
	}, nil
}

func (s BackofficeService) Update(ctx context.Context, req sellparam.UpdateSellRequest) (sellparam.UpdateSellResponse, error) {
	const op = "sellservice.Update"
	r, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return sellparam.UpdateSellResponse{}, richerror.New(op).WithWarpError(err)
	}
	r.Status = req.Status
	r.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, r)
	if err != nil {
		return sellparam.UpdateSellResponse{}, richerror.New(op).WithWarpError(err)
	}
	return sellparam.UpdateSellResponse{}, nil
}

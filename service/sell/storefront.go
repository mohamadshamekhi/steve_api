package sellservice

import (
	"context"
	"sync"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
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

type StorefrontRepo interface {
	Add(ctx context.Context, sell entity.Sell) error
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

func (s StorefrontService) Create(ctx context.Context, req sellparam.AddSellRequest) (sellparam.AddSellResponse, error) {
	const op = "sellservice.Create"
	sell := entity.Sell{
		UserID:            req.UserID,
		BrandID:           req.BrandID,
		ProductID:         req.ProductID,
		Color:             req.Color,
		BatteryHealth:     req.BatteryHealth,
		RegisterStatus:    req.RegisterStatus,
		PhysicalCondition: req.PhysicalCondition,
		Description:       req.Description,
		Status:            entity.SellStatusNew,
	}
	sell.BeforeCreate()
	err := s.Repo.Add(ctx, sell)
	if err != nil {
		return sellparam.AddSellResponse{}, richerror.New(op).WithWarpError(err)
	}

	return sellparam.AddSellResponse{}, nil
}

func (s StorefrontService) GetMySell(ctx context.Context, req sellparam.GetAllMySellParamRequest) (sellparam.GetAllMySellParamResponse, error) {
	const op = "sellservice.GetMySell"
	resp, err := s.Repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return sellparam.GetAllMySellParamResponse{}, richerror.New(op).WithWarpError(err)
	}

	var (
		userIDs    []string
		brandIDs   []string
		productIDs []string
	)
	for _, sell := range resp.Sells {
		userIDs = append(userIDs, sell.UserID)
		brandIDs = append(brandIDs, sell.BrandID)
		productIDs = append(productIDs, sell.ProductID)
	}

	var (
		uResp userparam.GetUsersByIDsResponse
		bResp brandparam.GetBrandsByIDsResponse
		pResp productparam.GetProductsByIDsResponse
	)

	var wg sync.WaitGroup
	wg.Add(3)

	// Users
	go func() {
		defer wg.Done()
		res, err := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
		if err != nil {
			logger.L().Error("can't fetch users - err:", err)
			return
		}
		uResp = res
	}()

	// Brands
	go func() {
		defer wg.Done()
		res, err := s.bProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: brandIDs})
		if err != nil {
			logger.L().Error("can't fetch brands - err:", err)
			return
		}
		bResp = res
	}()

	// Products
	go func() {
		defer wg.Done()
		res, err := s.pProvider.GetProductsByIDs(ctx, productparam.GetProductsByIDsRequest{IDs: productIDs})
		if err != nil {
			logger.L().Error("can't fetch products - err:", err)
			return
		}
		pResp = res
	}()

	wg.Wait()

	// ساخت Map برای دسترسی سریع
	userMap := make(map[string]userparam.UserInfo, len(uResp.Users))
	for _, u := range uResp.Users {
		userMap[u.ID] = u
	}

	brandMap := make(map[string]entity.Brand, len(bResp.Brands))
	for id, b := range bResp.Brands {
		brandMap[id] = b
	}

	productMap := make(map[string]productparam.ProductList, len(pResp.Products))
	for id, p := range pResp.Products {
		productMap[id] = p
	}

	// ساختن خروجی Sells
	sellInfos := make([]sellparam.SellInfo, 0, len(resp.Sells))
	for _, r := range resp.Sells {
		user := userMap[r.UserID]
		brand := brandMap[r.BrandID]
		product := productMap[r.ProductID]

		sellInfos = append(sellInfos, sellparam.SellInfo{
			ID: r.ID,
			User: struct {
				Name        string `json:"name"`
				PhoneNumber string `json:"phone_number"`
			}{
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			BrandName:         brand.Name,
			ProductName:       product.Name,
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

	return sellparam.GetAllMySellParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Sells: sellInfos,
	}, nil
}

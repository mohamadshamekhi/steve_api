package sellservice

import (
	"context"
	"sync"
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
	wg := sync.WaitGroup{}
	var mu sync.Mutex

	wg.Add(3)
	user := userparam.UserInfo{}
	product := productparam.ProductList{}
	brand := entity.Brand{}

	go func() {
		defer wg.Done()

		bResp, uErr := s.bProvider.GetBrandsByIDs(ctx, brandparam.GetBrandsByIDsRequest{IDs: []string{r.BrandID}})
		if uErr != nil {
			logger.L().Error("Failed to get repair brand info", uErr)
			return
		}
		if b, ok := bResp.Brands[r.BrandID]; ok {
			mu.Lock()
			brand = b
			mu.Unlock()
		}

	}()
	go func() {
		defer wg.Done()

		pResp, pErr := s.pProvider.GetProductsByIDs(ctx, productparam.GetProductsByIDsRequest{IDs: []string{r.ProductID}})
		if pErr != nil {
			logger.L().Error("Failed to get repair products", pErr)
			return
		}
		if p, ok := pResp.Products[r.ProductID]; ok {
			mu.Lock()
			product = p
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()

		uResp, uErr := s.uProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: []string{r.UserID}})
		if uErr != nil {
			logger.L().Error("Failed to get repair users", uErr)
			return
		}
		for _, u := range uResp.Users {
			if u.ID == r.UserID {
				mu.Lock()
				user = u
				mu.Unlock()
			}
		}
	}()

	wg.Wait()

	return sellparam.GetSellResponse{Sell: sellparam.SellInfo{
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
	}}, nil
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

func (s BackofficeService) GetAll(ctx context.Context, req sellparam.GetAllSellParamRequest) (sellparam.GetAllSellParamResponse, error) {
	const op = "sellservice.GetAll"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return sellparam.GetAllSellParamResponse{}, richerror.New(op).WithWarpError(err)
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

	return sellparam.GetAllSellParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Sells: sellInfos,
	}, nil
}

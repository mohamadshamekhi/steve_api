package repairservice

import (
	"context"
	"sync"
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

	return repairparam.GetRepairResponse{RepairInfo: repairparam.RepairInfo{
		ID: r.ID,
		User: struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		BrandName:   brand.Name,
		ProductName: product.Name,
		Galleries:   r.Galleries,
		Description: r.Description,
		Status:      r.Status,
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

	var (
		userIDs    []string
		brandIDs   []string
		productIDs []string
	)
	for _, sell := range resp.Repairs {
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

	repairInfo := make([]repairparam.RepairInfo, 0, len(resp.Repairs))
	for _, r := range resp.Repairs {
		user := userMap[r.UserID]
		brand := brandMap[r.BrandID]
		product := productMap[r.ProductID]

		repairInfo = append(repairInfo, repairparam.RepairInfo{
			ID: r.ID,
			User: struct {
				Name        string `json:"name"`
				PhoneNumber string `json:"phone_number"`
			}{
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			BrandName:   brand.Name,
			ProductName: product.Name,
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
		Repairs: repairInfo,
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

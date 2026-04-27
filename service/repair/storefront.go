package repairservice

import (
	"context"
	"sync"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	repairparam "github.com/abolfazlnorzad/steve/param/repair"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	persian "github.com/yaa110/go-persian-calendar"
)

type StorefrontService struct {
	Repo                 StorefrontRepo
	uProvider            UserProvider
	bProvider            BrandProvider
	pProvider            ProductProvider
	backofficeProvider   BackofficeProvider
	melliPayamakProvider MelliPayamakProvider
}

type StorefrontRepo interface {
	Add(ctx context.Context, repair entity.Repair) error
	All(ctx context.Context, filter commonparam.FilterRequest, pagination commonparam.PaginationRequest) (ListResult, error)
	GetUserRepairRequestStatics(ctx context.Context, userID string) (uint64, error)
}

func NewStorefrontService(r StorefrontRepo, uProvider UserProvider,
	bProvider BrandProvider,
	pProvider ProductProvider,
	backofficeProvider BackofficeProvider,
	melliPayamakProvider MelliPayamakProvider,
) StorefrontService {
	return StorefrontService{
		Repo:                 r,
		uProvider:            uProvider,
		bProvider:            bProvider,
		pProvider:            pProvider,
		backofficeProvider:   backofficeProvider,
		melliPayamakProvider: melliPayamakProvider,
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

	loc, _ := time.LoadLocation("Asia/Tehran")
	now := time.Now().In(loc)
	pt := persian.New(now)
	date := pt.Format("yyyy/MM/dd")
	hour := pt.Format("HH:mm")

	resp, bErr := s.backofficeProvider.GetSuperUserBackOffices(ctx, backofficeparam.GetSuperUserBackofficeRequest{})
	if bErr != nil {
		logger.L().Error("Failed to get super user backoffice", "err", err)
	}

	user, uErr := s.uProvider.GetUserByID(ctx, userparam.GetUsersRequest{ID: req.UserID})
	if uErr != nil {
		logger.L().Error("Failed to get user", "err", uErr)
	}

	for _, backoffice := range resp.Backoffice {
		go func(date string, hour string) {
			sCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if mErr := s.melliPayamakProvider.SendPatternSMS(sCtx, backoffice.PhoneNumber, []string{
				backoffice.Name,
				user.User.Name,
				user.User.PhoneNumber,
				"کالایی برای تعمیرات",
				"تلفن",
				user.User.PhoneNumber,
				date, hour,
			}, mellipayamak.BackofficeKey); mErr != nil {
				logger.L().Error("Failed to send sms backoffice", "err", mErr)
			}
		}(date, hour)
	}

	return repairparam.AddRepairResponse{}, nil
}

func (s StorefrontService) GetAllRepair(ctx context.Context, req repairparam.GetAllMyRepairParamRequest) (repairparam.GetAllMyRepairParamResponse, error) {
	const op = "StorefrontService.GetAllRepair"
	resp, err := s.Repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return repairparam.GetAllMyRepairParamResponse{}, richerror.New(op).WithWarpError(err)
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
			CreatedAt:   r.CreatedAt,
			Galleries:   r.Galleries,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return repairparam.GetAllMyRepairParamResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(resp.TotalCount),
		},
		Repairs: repairInfo,
	}, nil
}

func (s StorefrontService) GetUserRepairStatics(ctx context.Context, req repairparam.GetUserReairStaticRequest) (repairparam.GetUserReairStaticResponse, error) {
	const op = "sellhandler.GetUserRepairStatics"
	c, err := s.Repo.GetUserRepairRequestStatics(ctx, req.UserID)
	if err != nil {
		return repairparam.GetUserReairStaticResponse{}, err
	}
	return repairparam.GetUserReairStaticResponse{
		Count: c,
	}, nil
}

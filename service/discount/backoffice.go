package discountservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/gommon/log"
)

type BackofficeService struct {
	repo                DiscountReadWriter
	cartInfoProvider    CartInfoProvider
	productInfoProvider ProductInfoProvider
}

type CartInfoProvider interface {
	GetEnrichedCart(ctx context.Context, cart *entity.Cart) (EnrichedCart, error)
}
type ProductInfoProvider interface {
	GetProductCategoriesForCartItems(ctx context.Context, req productparam.GetProductCategoriesForCartItemsRequest) (productparam.GetProductCategoriesForCartItemsResponse, error)
}

type EnrichedCart struct {
	*entity.Cart
	ProductIDs  []string
	CategoryIDs []string
	Subtotal    int64
}

func NewBackofficeService(repo DiscountReadWriter, cartInfoProvider CartInfoProvider,
	productInfoProvider ProductInfoProvider) BackofficeService {
	return BackofficeService{repo: repo, cartInfoProvider: cartInfoProvider, productInfoProvider: productInfoProvider}
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
		return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithWarpError(err)
	}

	return discountparam.UpdateDiscountResponse{}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req discountparam.GetAllDiscountsParamRequest) (discountparam.GetAllDiscountsParamResponse, error) {
	const op = "backofficeservice.GetAll"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return discountparam.GetAllDiscountsParamResponse{}, richerror.New(op).WithWarpError(err)
	}

	return discountparam.GetAllDiscountsParamResponse{Discounts: resp.Discounts, Pagination: struct {
		PageSize   uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
		Total      uint `json:"total"`
	}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)}}, nil
}

func (s BackofficeService) GetSingle(ctx context.Context, req discountparam.GetSingleDiscountRequest) (discountparam.GetSingleDiscountResponse, error) {
	const op = "backofficeservice.GetSingle"
	discount, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return discountparam.GetSingleDiscountResponse{}, richerror.New(op).WithWarpError(err)
	}

	return discountparam.GetSingleDiscountResponse{
		Discount: discount,
	}, nil

}

func (s BackofficeService) GetByCode(ctx context.Context, req discountparam.GetSingleDiscountByCodeRequest) (discountparam.GetSingleDiscountByCodeResponse, error) {
	const op = "backofficeservice.GetByCode"
	discount, err := s.repo.FindByCode(ctx, req.Code)
	if err != nil {
		return discountparam.GetSingleDiscountByCodeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return discountparam.GetSingleDiscountByCodeResponse{
		Discount: discount,
	}, nil
}

func (s BackofficeService) GetValidDiscountsForCart(ctx context.Context, req discountparam.GetValidDiscountsForSpecificCartRequest) (discountparam.GetValidDiscountsForSpecificCartResponse, error) {
	const op = "backofficeservice.GetValidDiscountsForCart"
	enrichedCart, err := s.cartInfoProvider.GetEnrichedCart(ctx, &req.Cart)
	if err != nil {
		return discountparam.GetValidDiscountsForSpecificCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	potentialDiscounts, err := s.repo.FindApplicableDiscounts(ctx, req.DiscountCode)
	if err != nil {
		return discountparam.GetValidDiscountsForSpecificCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	validDiscounts := make([]entity.Discount, 0)
	for _, discount := range potentialDiscounts {
		if s.isDiscountValid(ctx, discount, enrichedCart) {
			validDiscounts = append(validDiscounts, discount)
		}
	}

	return discountparam.GetValidDiscountsForSpecificCartResponse{Discounts: validDiscounts}, nil
}

func (s BackofficeService) isDiscountValid(ctx context.Context, d entity.Discount, cart EnrichedCart) bool {
	now := time.Now()

	if !d.IsActive || now.Before(d.StartsAt) || now.After(d.ExpiresAt) {
		return false
	}
	if d.TotalUsageLimit > 0 && d.UsageCount >= d.TotalUsageLimit {
		return false
	}
	if cart.Subtotal < int64(d.MinimumPurchaseAmount) {
		return false
	}

	if len(d.ApplicableUserIDs) > 0 && !contains(d.ApplicableUserIDs, cart.UserID) {
		return false
	}

	if len(d.ApplicableProductIDs) > 0 && !intersects(d.ApplicableProductIDs, cart.ProductIDs) {
		return false
	}

	if len(d.ApplicableCategoryIDs) > 0 && !intersects(d.ApplicableCategoryIDs, cart.CategoryIDs) {
		return false
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func intersects(slice1, slice2 []string) bool {
	set := make(map[string]struct{})
	for _, s := range slice1 {
		set[s] = struct{}{}
	}
	for _, s := range slice2 {
		if _, ok := set[s]; ok {
			return true
		}
	}
	return false
}

func (s BackofficeService) CalculateDiscountAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error) {
	eligibleAmount, err := s.calculateEligibleAmount(ctx, d, cart)
	if err != nil {
		return 0, err
	}

	var discountAmount int64 = 0
	switch d.Type {
	case entity.PercentageDiscount:
		if eligibleAmount > 0 && d.Value > 0 {
			discountAmount = (eligibleAmount * int64(d.Value)) / 100
		}
	case entity.FixedAmountDiscount:
		discountAmount = int64(d.Value)
	}

	if discountAmount > eligibleAmount {
		return eligibleAmount, nil
	}

	return discountAmount, nil
}

func (s BackofficeService) calculateEligibleAmount(ctx context.Context, d entity.Discount, cart entity.Cart) (int64, error) {
	isGeneralDiscount := len(d.ApplicableProductIDs) == 0 && len(d.ApplicableCategoryIDs) == 0
	if isGeneralDiscount {
		return int64(cart.CalculateSubtotal()), nil
	}

	var eligibleAmount int64 = 0

	productToCategoriesMap, err := s.productInfoProvider.GetProductCategoriesForCartItems(ctx, productparam.GetProductCategoriesForCartItemsRequest{Items: cart.Items})
	if err != nil {
		return 0, err
	}

	productIDSet := make(map[string]struct{})
	for _, id := range d.ApplicableProductIDs {
		productIDSet[id] = struct{}{}
	}
	categoryIDSet := make(map[string]struct{})
	for _, id := range d.ApplicableCategoryIDs {
		categoryIDSet[id] = struct{}{}
	}

	for _, item := range cart.Items {
		isEligible := false

		if _, ok := productIDSet[item.ProductID]; ok {
			isEligible = true
		}

		if !isEligible {
			if categoryIDs, ok := productToCategoriesMap.CategoryIDs[item.ProductID]; ok {
				for _, catID := range categoryIDs {
					if _, catOK := categoryIDSet[catID]; catOK {
						isEligible = true
						break
					}
				}
			}
		}

		if isEligible {
			eligibleAmount += item.PriceAtTime * int64(item.Quantity)
		}
	}

	return eligibleAmount, nil
}

func (s BackofficeService) UpdateDiscountUsageAfterRegisterOrder(ctx context.Context, req discountparam.UpdateDiscountUsageAfterRegisterOrderRequest) (discountparam.UpdateDiscountUsageAfterRegisterOrderResponse, error) {
	const op = "discount.UpdateDiscountUsageAfterRegisterOrder"
	for _, d := range req.AppliedDiscount {
		if err := s.repo.UpdateUsageCount(ctx, d.DiscountID); err != nil {
			log.Error(ctx, op, err, "failed to update discount usage - discount id", d.DiscountID)
		}
	}
	return discountparam.UpdateDiscountUsageAfterRegisterOrderResponse{}, nil
}

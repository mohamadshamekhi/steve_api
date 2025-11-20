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
	FindProductWithSpecialCategory(ctx context.Context, req productparam.GetProductWithSpecialCategoryRequest) (productparam.GetProductWithSpecialCategoryResponse, error)
	GetPricesOfProducts(ctx context.Context, req productparam.GetPriceOfProductsRequest) (productparam.GetPriceOfProductsResponse, error)
	GetProductIDsByCategories(ctx context.Context, req productparam.GetProductIDsByCategoriesRequest) (productparam.GetProductIDsByCategoriesResponse, error)
}

type EnrichedCart struct {
	*entity.Cart
	ProductIDs  map[string][]string
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
		IsSpecial:             req.IsSpecial,
		ApplicableTo:          entity.ApplicableTo(req.ApplicableTo),
		ApplicableProductIDs:  req.ApplicableProductIDs,
		ApplicableCategoryIDs: req.ApplicableCategoryIDs,
		ApplicableUserIDs:     req.ApplicableUserIDs,
		IsActive:              req.IsActive,
	}
	d.BeforeCreate()

	// check for applicable_to
	if _, err := s.validateApplicableTo(d, req); err != nil {
		return discountparam.AddDiscountResponse{}, richerror.New(op).WithWarpError(err)
	}

	// check for special discount
	if req.IsSpecial {
		r, e := s.repo.ExistSpecialDiscount(ctx, req.StartsAt, req.ExpiresAt)
		if e != nil {
			return discountparam.AddDiscountResponse{}, richerror.New(op).WithWarpError(e)
		}
		if r {
			return discountparam.AddDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.special_discount_is_exist_in_this_time")
		}
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return discountparam.AddDiscountResponse{}, err
	}
	return discountparam.AddDiscountResponse{}, nil
}

func (s BackofficeService) validateApplicableTo(d entity.Discount, req discountparam.AddDiscountRequest) (discountparam.AddDiscountResponse, error) {
	const op = "backoffice.validateApplicableTo"
	switch d.ApplicableTo {
	case entity.ApplicableToCategories:
		if req.ApplicableCategoryIDs == nil {
			return discountparam.AddDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_category_required")
		} else {
			// todo - check category is exist
		}
	case entity.ApplicableToUsers:
		if req.ApplicableUserIDs == nil {
			return discountparam.AddDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_user_required")
		} else {
			// todo - check user is exist
		}
	case entity.ApplicableToProducts:
		if req.ApplicableProductIDs == nil {
			return discountparam.AddDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_product_required")
		} else {
			// todo - check product & variant is exist
		}

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
	if req.IsSpecial != nil {
		if *req.IsSpecial {
			var ex time.Time
			if req.ExpiresAt != nil {
				ex = *req.ExpiresAt
			} else {
				ex = discount.ExpiresAt
			}
			r, e := s.repo.ExistSpecialDiscount(ctx, discount.StartsAt, ex)
			if e != nil {
				return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithWarpError(e)
			}
			if r {
				return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.special_discount_is_exist_in_this_time")
			}
		}
		discount.IsSpecial = *req.IsSpecial
	}
	if req.ApplicableTo != nil {
		discount.ApplicableTo = entity.ApplicableTo(*req.ApplicableTo)
		switch *req.ApplicableTo {
		case string(entity.ApplicableToCategories):
			if req.ApplicableCategoryIDs == nil {
				return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_category_required")
			} else {
				discount.ApplicableCategoryIDs = *req.ApplicableCategoryIDs
				// todo - check category is exist
			}
		case string(entity.ApplicableToUsers):
			if req.ApplicableUserIDs == nil {
				return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_user_required")
			} else {
				discount.ApplicableUserIDs = *req.ApplicableUserIDs
				// todo - check user is exist
			}
		case string(entity.ApplicableToProducts):
			if req.ApplicableProductIDs == nil {
				return discountparam.UpdateDiscountResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("discount.applicable_to_product_required")
			} else {
				discount.ApplicableProductIDs = *req.ApplicableProductIDs
				// todo - check product & variant is exist
			}

		}
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

	switch d.ApplicableTo {
	case entity.ApplicableToUsers:
		if len(d.ApplicableUserIDs) > 0 && !contains(d.ApplicableUserIDs, cart.UserID) {
			return false
		}
	case entity.ApplicableToCategories:
		if len(d.ApplicableCategoryIDs) > 0 && !intersects(d.ApplicableCategoryIDs, cart.CategoryIDs) {
			return false
		}
	case entity.ApplicableToProducts:
		if len(d.ApplicableProductIDs) > 0 && !HasCommonVariant(d.ApplicableProductIDs, cart.ProductIDs) {
			return false
		}

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

func HasCommonVariant(a []entity.SpecialDiscountItem, b map[string][]string) bool {
	for _, item := range a {
		// بررسی وجود ProductID در map
		if variants, exists := b[item.ProductID]; exists {
			// اگر در تخفیف، VariantID وجود داره (nil نیست)
			if item.VariantID != nil {
				for _, v := range variants {
					if v == *item.VariantID {
						return true // اشتراک پیدا شد
					}
				}
			} else {
				// اگر VariantID مشخص نشده، یعنی تخفیف برای کل محصوله
				return true
			}
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

func (s BackofficeService) calculateEligibleAmount(
	ctx context.Context,
	d entity.Discount,
	cart entity.Cart,
) (int64, error) {

	switch d.ApplicableTo {
	case entity.ApplicableToProducts:
		// SpecialDiscountItems: ProductID + (VariantID اختیاری = همه ورینت‌ها)
		productsAllVariants := make(map[string]struct{})
		productVariantSet := make(map[string]map[string]struct{})
		for _, sdi := range d.ApplicableProductIDs {
			if sdi.VariantID == nil {
				productsAllVariants[sdi.ProductID] = struct{}{}
				continue
			}
			if productVariantSet[sdi.ProductID] == nil {
				productVariantSet[sdi.ProductID] = make(map[string]struct{})
			}
			productVariantSet[sdi.ProductID][*sdi.VariantID] = struct{}{}
		}

		var eligible int64
		for _, item := range cart.Items {
			if _, ok := productsAllVariants[item.ProductID]; ok {
				eligible += item.PriceAtTime * int64(item.Quantity)
				continue
			}
			if vs, ok := productVariantSet[item.ProductID]; ok {
				if _, ok := vs[item.VariantID]; ok {
					eligible += item.PriceAtTime * int64(item.Quantity)
				}
			}
		}
		return eligible, nil

	case entity.ApplicableToCategories:
		if len(d.ApplicableCategoryIDs) == 0 {
			return 0, nil
		}
		categoryIDSet := make(map[string]struct{}, len(d.ApplicableCategoryIDs))
		for _, id := range d.ApplicableCategoryIDs {
			categoryIDSet[id] = struct{}{}
		}
		// فقط وقتی لازم است دسته‌ها را می‌گیریم
		resp, err := s.productInfoProvider.GetProductCategoriesForCartItems(
			ctx,
			productparam.GetProductCategoriesForCartItemsRequest{Items: cart.Items},
		)
		if err != nil {
			return 0, err
		}

		var eligible int64
		for _, item := range cart.Items {
			if cats, ok := resp.CategoryIDs[item.ProductID]; ok {
				for _, c := range cats {
					if _, hit := categoryIDSet[c]; hit {
						eligible += item.PriceAtTime * int64(item.Quantity)
						break
					}
				}
			}
		}
		return eligible, nil

	case entity.ApplicableToUsers:
		// این مسیر معمولاً برای مجوزدهی/واجدمعرفیته نه مبلغ‌سنجی.
		// اگر لازم داری مبلغ فقط وقتی کاربر در لیسته محاسبه بشه:
		if len(d.ApplicableUserIDs) == 0 {
			return 0, nil
		}
		if !contains(d.ApplicableUserIDs, cart.UserID) {
			return 0, nil
		}
		return int64(cart.CalculateSubtotal()), nil

	default:
		// ناشناخته: هیچ‌چیزی محاسبه نکن
		return 0, nil
	}
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

func (s BackofficeService) GetProductIDsThatHaveActiveAutoDiscount(
	ctx context.Context,
	_ discountparam.GetProductIDsThatHaveActiveAutoDiscountRequest,
) (discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse, error) {

	const op = "discountservice.GetProductIDsThatHaveActiveAutoDiscount"

	// مرحله ۱: دریافت تمام تخفیف‌های خودکار فعال
	autoDiscounts, err := s.repo.GetActiveAutoDiscount(ctx)
	if err != nil {
		return discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse{}, richerror.New(op).WithWarpError(err)
	}

	// مجموعه‌ای از شناسه‌ها برای دریافت اطلاعات به صورت یکجا (Batch)
	productIDsForPrice := make(map[string]struct{})
	variantIDsForPrice := make(map[string]struct{})
	categoryIDs := make(map[string]struct{})

	for _, d := range autoDiscounts {
		switch d.ApplicableTo {
		case entity.ApplicableToProducts:
			for _, p := range d.ApplicableProductIDs {
				if p.ProductID != "" {
					productIDsForPrice[p.ProductID] = struct{}{}
					if p.VariantID != nil && *p.VariantID != "" {
						variantIDsForPrice[*p.VariantID] = struct{}{}
					}
				}
			}
		case entity.ApplicableToCategories:
			for _, cid := range d.ApplicableCategoryIDs {
				if cid != "" {
					categoryIDs[cid] = struct{}{}
				}
			}
		}
	}

	// تبدیل map به slice برای ارسال به سرویس‌ها
	catIDSlice := make([]string, 0, len(categoryIDs))
	for id := range categoryIDs {
		catIDSlice = append(catIDSlice, id)
	}

	// مرحله ۲: دریافت محصولات مربوط به دسته‌بندی‌ها
	categoryToProductIDs := make(map[string][]string) // map[categoryID][]productID
	if len(catIDSlice) > 0 {
		prodsByCat, err := s.productInfoProvider.GetProductIDsByCategories(ctx, productparam.GetProductIDsByCategoriesRequest{CatIDs: catIDSlice})
		if err != nil {
			return discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse{}, richerror.New(op).WithWarpError(err)
		}
		for catID, prodIDs := range prodsByCat.ProdsByCat {
			categoryToProductIDs[catID] = prodIDs
			for _, pid := range prodIDs {
				productIDsForPrice[pid] = struct{}{}
			}
		}
	}

	// تبدیل map به slice برای ارسال به سرویس قیمت
	prodIDSlice := make([]string, 0, len(productIDsForPrice))
	for id := range productIDsForPrice {
		prodIDSlice = append(prodIDSlice, id)
	}
	variantIDSlice := make([]string, 0, len(variantIDsForPrice))
	for id := range variantIDsForPrice {
		variantIDSlice = append(variantIDSlice, id)
	}

	// مرحله ۳: دریافت قیمت تمام محصولات و ورینت‌های مورد نیاز به صورت یکجا
	prices, pErr := s.productInfoProvider.GetPricesOfProducts(ctx, productparam.GetPriceOfProductsRequest{ProductIDs: prodIDSlice, VariantIDs: variantIDSlice})

	if pErr != nil {
		return discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse{}, richerror.New(op).WithWarpError(err)
	}

	// مرحله ۴: پردازش تخفیف‌ها و انتخاب بهترین تخفیف برای هر ترکیب (محصول-تنوع)
	type BestDiscountInfo struct {
		Discount productparam.ProductDiscount
		Value    int // مقدار ریالی تخفیف برای مقایسه
	}
	// ساختار کلیدی: map[ProductID]map[VariantID]BestDiscountInfo
	bestDiscountsPerVariant := make(map[string]map[string]BestDiscountInfo)

	// تابع کمکی برای به‌روزرسانی بهترین تخفیف
	updateBestDiscount := func(pID, vID string, price float64, d entity.Discount) {
		var calculatedValue int
		if d.Type == entity.PercentageDiscount {
			calculatedValue = int(price * float64(d.Value) / 100)
		} else { // FixedAmountDiscount
			calculatedValue = int(d.Value)
		}

		if _, ok := bestDiscountsPerVariant[pID]; !ok {
			bestDiscountsPerVariant[pID] = make(map[string]BestDiscountInfo)
		}

		if existing, exists := bestDiscountsPerVariant[pID][vID]; !exists || calculatedValue > existing.Value {

			bestDiscountsPerVariant[pID][vID] = BestDiscountInfo{
				Value: calculatedValue,
				Discount: productparam.ProductDiscount{
					DiscountID:    d.ID,
					DiscountType:  string(d.Type),
					DiscountValue: d.Value,
					StartAt:       d.StartsAt,
					ExpiresAt:     d.ExpiresAt,
					IsActive:      d.IsActive,
					IsSpecial:     d.IsSpecial,
					ApplicableTo:  d.ApplicableTo,
					ProductID:     pID,
					VariantID:     vID,
				},
			}
		}
	}

	for _, d := range autoDiscounts {
		switch d.ApplicableTo {
		case entity.ApplicableToProducts:
			for _, p := range d.ApplicableProductIDs {
				if p.ProductID == "" {
					continue
				}
				var t string
				var v uint
				if p.DiscountType != "" {
					t = string(p.DiscountType)
				} else {
					t = string(d.Type)
				}

				if p.DiscountValue != 0 {
					v = p.DiscountValue
				} else {
					v = d.Value
				}
				if p.VariantID != nil && *p.VariantID != "" {
					variantID := *p.VariantID
					if variantPrices, ok := prices.Price[p.ProductID]; ok {
						if price, ok := variantPrices[variantID]; ok {
							updateBestDiscount(p.ProductID, variantID, price, entity.Discount{
								ID:           d.ID,
								Type:         entity.DiscountType(t),
								Value:        v,
								StartsAt:     d.StartsAt,
								ExpiresAt:    d.ExpiresAt,
								IsActive:     d.IsActive,
								IsSpecial:    d.IsSpecial,
								ApplicableTo: d.ApplicableTo,
							})
						}
					}
				} else {
					if variantPrices, ok := prices.Price[p.ProductID]; ok {
						for variantID, price := range variantPrices {
							updateBestDiscount(p.ProductID, variantID, price, d)
						}
					}
				}
			}

		case entity.ApplicableToCategories:
			for _, cid := range d.ApplicableCategoryIDs {
				productIDs, ok := categoryToProductIDs[cid]
				if !ok {
					continue
				}
				for _, pid := range productIDs {
					if variantPrices, ok := prices.Price[pid]; ok {
						for variantID, price := range variantPrices {
							updateBestDiscount(pid, variantID, price, d)
						}
					}
				}
			}
		}
	}

	// ========================>> مرحله ۵: بازنویسی شده <<========================
	// ساختار نهایی برای نگهداری تخفیف‌های هر ورینت به صورت مجزا
	// کلید: ProductID، مقدار: لیستی از تخفیف‌های ورینت‌های آن محصول
	finalDiscountsMap := make(map[string][]productparam.ProductDiscount)

	for productID, variantDiscounts := range bestDiscountsPerVariant {
		// برای هر محصول، یک لیست از تخفیف‌های ورینت‌هایش ایجاد می‌کنیم
		discountsForThisProduct := make([]productparam.ProductDiscount, 0, len(variantDiscounts))

		// بهترین تخفیف محاسبه شده برای هر ورینت را به لیست اضافه می‌کنیم
		for _, discountInfo := range variantDiscounts {
			discountsForThisProduct = append(discountsForThisProduct, discountInfo.Discount)
		}

		// اگر حداقل یک تخفیف برای ورینت‌های این محصول وجود داشت، آن را به نقشه نهایی اضافه می‌کنیم
		if len(discountsForThisProduct) > 0 {
			finalDiscountsMap[productID] = discountsForThisProduct
		}
	}
	// =======================================================================

	// نکته مهم: نوع فیلد DiscountMap در استراکت GetProductIDsThatHaveActiveAutoDiscountResponse
	// باید به map[string][]productparam.ProductDiscount تغییر کند.
	return discountparam.GetProductIDsThatHaveActiveAutoDiscountResponse{
		DiscountMap: finalDiscountsMap,
	}, nil
}

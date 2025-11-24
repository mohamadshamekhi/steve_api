package productservice

import (
	"context"

	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s StorefrontService) GetProductsWithSpecialDiscount(
	ctx context.Context,
	_ productparam.GetProductsWithSpecialDiscountRequest,
) (productparam.GetProductsWithSpecialDiscountResponse, error) {
	const op = "productservice.GetProductsWithSpecialDiscount"

	// ۱) گرفتن تخفیف شگفت‌انگیز فعال
	specialDiscountResp, err := s.discountProvider.GetSpecialDiscount(
		ctx,
		discountparam.GetSpecialDiscountRequest{},
	)
	if err != nil {
		return productparam.GetProductsWithSpecialDiscountResponse{},
			richerror.New(op).WithWarpError(err)
	}

	d := specialDiscountResp.Discount
	if !d.IsSpecial || !d.IsActive {
		// اگر چیزی نبود، خروجی خالی
		return productparam.GetProductsWithSpecialDiscountResponse{
			Products: []productparam.StorefrontProductList{},
		}, nil
	}

	// ۲) ساخت map برای جلوگیری از تکرار محصول
	productMap := make(map[string]*productparam.StorefrontProductList)

	for _, item := range d.ApplicableProductIDs {
		if item.ProductID == "" {
			continue
		}

		// اگر محصول قبلاً نخوندیم، از repo بگیر و داخل map ذخیره کن
		p, ok := productMap[item.ProductID]
		if !ok {
			prod, err := s.repo.FindByID(ctx, item.ProductID)
			if err != nil {
				return productparam.GetProductsWithSpecialDiscountResponse{},
					richerror.New(op).WithWarpError(err)
			}

			p = &productparam.StorefrontProductList{
				ID:           prod.ID,
				Name:         prod.Name,
				ENTitle:      prod.ENTitle,
				Slug:         prod.Slug,
				Gallery:      prod.Gallery,
				Rating:       prod.Rating,
				IsPublished:  prod.IsPublished,
				IsSecondHand: prod.IsSecondHand,
				Variant:      prod.Variants,
				Discount:     []productparam.ProductDiscount{},
			}
			productMap[prod.ID] = p
		}

		// ۳) تعیین نوع و مقدار تخفیف (overridden یا عمومی)
		discType := d.Type
		if item.DiscountType != "" {
			discType = item.DiscountType
		}

		discValue := d.Value
		if item.DiscountValue != 0 {
			discValue = item.DiscountValue
		}

		// ۴) اگر VariantID مشخص است → فقط روی همان ورینت
		if item.VariantID != nil && *item.VariantID != "" {
			p.Discount = append(p.Discount, productparam.ProductDiscount{
				DiscountID:    d.ID,
				DiscountType:  string(discType),
				DiscountValue: discValue,
				StartAt:       d.StartsAt,
				ExpiresAt:     d.ExpiresAt,
				IsActive:      d.IsActive,
				IsSpecial:     d.IsSpecial,
				ApplicableTo:  d.ApplicableTo,
				VariantID:     *item.VariantID,
				ProductID:     item.ProductID,
			})
		} else {
			// ۵) اگر VariantID خالی بود → تخفیف روی تمام ورینت‌های محصول
			for _, v := range p.Variant {
				p.Discount = append(p.Discount, productparam.ProductDiscount{
					DiscountID:    d.ID,
					DiscountType:  string(discType),
					DiscountValue: discValue,
					StartAt:       d.StartsAt,
					ExpiresAt:     d.ExpiresAt,
					IsActive:      d.IsActive,
					IsSpecial:     d.IsSpecial,
					ApplicableTo:  d.ApplicableTo,
					VariantID:     v.ID,
					ProductID:     item.ProductID,
				})
			}
		}
	}

	// ۶) تبدیل map به slice برای خروجی
	res := make([]productparam.StorefrontProductList, 0, len(productMap))
	for _, p := range productMap {
		res = append(res, *p)
	}

	return productparam.GetProductsWithSpecialDiscountResponse{
		Products: res,
	}, nil
}

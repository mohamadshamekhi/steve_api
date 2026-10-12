package cartservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
)

type StoreFrontService struct {
	repo          CartReadWriter
	pRepo         PRepository
	uRepo         URepository
	orderProvider orderservice.StoreFrontService
}

func NewStorefrontService(repo CartReadWriter,
	pRepo PRepository, uRepo URepository) StoreFrontService {
	return StoreFrontService{
		repo:  repo,
		pRepo: pRepo,
		uRepo: uRepo,
	}
}

func (s *StoreFrontService) SetOrderProvider(orderProvider orderservice.StoreFrontService) {
	s.orderProvider = orderProvider
}

func (s *StoreFrontService) AddItemToCart(ctx context.Context, req cartparam.AddItemToCartRequest) (cartparam.AddItemToCartResponse, error) {
	const op = "cartservice.AddItemToCart"
	p, _ := s.pRepo.FindByID(ctx, req.ProductID)
	userID, _ := ctx.Value("user_id").(string)
	var price uint64
	var variants []entity.VariantAttribute
	for _, variant := range p.Variants {
		if variant.ID == req.VariantID {
			price = variant.Price
			variants = variant.VariantAttribute
		}
	}
	item := entity.CartItem{
		ProductID:         req.ProductID,
		VariantID:         req.VariantID,
		ProductName:       p.Name,
		VariantAttributes: variants,
		PriceAtTime:       int64(price),
	}
	if err := s.repo.AddToCart(ctx, userID, item); err != nil {
		return cartparam.AddItemToCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.AddItemToCartResponse{}, nil
}

func (s *StoreFrontService) RemoveFromCart(ctx context.Context, req cartparam.RemoveFromCartRequest) (cartparam.RemoveFromCartResponse, error) {
	const op = "cartservice.RemoveFromCart"
	userID, _ := ctx.Value("user_id").(string)
	if err := s.repo.RemoveFromCart(ctx, userID, req.ProductID, req.VariantID); err != nil {
		return cartparam.RemoveFromCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.RemoveFromCartResponse{}, nil
}

func (s *StoreFrontService) GetMyActiveCart(ctx context.Context, req cartparam.MyActiveCartRequest) (cartparam.MyActiveCartResponse, error) {
	const op = "cartservice.GetMyActiveCart"
	cart, err := s.repo.GetUserActiveCart(ctx, req.UserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return cartparam.MyActiveCartResponse{}, nil
		}
		return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	ids := []string{req.UserID}
	user, err := s.uRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	summary, err := s.orderProvider.FinalizeCartForCheckout(ctx, *cart, "")
	if err != nil {
		return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	var cartItems []cartparam.StoreFrontCartItemDetail
	for _, item := range cart.Items {
		detail := cartparam.StoreFrontCartItemDetail{
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			Summary:           summary,
		}
		if !item.IsCustom {
			p, pErr := s.pRepo.FindByID(ctx, item.ProductID)
			if pErr == nil {
				detail.ProductID = p.ID
				detail.Url = p.GetMainGallery().URL
				detail.RequiresFreightShipping = p.RequiresFreightShipping
			}
		}
		cartItems = append(cartItems, detail)
	}

	c := cartparam.StoreFrontCartDetail{
		ID:     cart.ID,
		Status: string(cart.Status),
		User: struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{
			ID:          cart.UserID,
			Name:        user[0].Name,
			PhoneNumber: user[0].PhoneNumber,
		},
		Items:     cartItems,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		Subtotal:  cart.CalculateSubtotal(),
	}
	return cartparam.MyActiveCartResponse{
		Cart: c,
	}, nil
}

func (s *StoreFrontService) GetEnrichedCart(
	ctx context.Context,
	cart *entity.Cart,
) (discountservice.EnrichedCart, error) {

	// productID -> []variantID
	productIDVariantIDs := make(map[string][]string)

	// برای اینکه دسته‌ها تکراری نشوند
	categoryIDSet := make(map[string]struct{})

	for _, item := range cart.Items {
		// Custom (hand-typed) lines have no catalog product behind them.
		if item.IsCustom {
			continue
		}

		// جمع کردن ورینت‌های هر محصول در map
		if item.ProductID != "" {
			variants := productIDVariantIDs[item.ProductID]

			// جلوگیری از تکرار VariantID داخل slice
			if item.VariantID != "" {
				already := false
				for _, v := range variants {
					if v == item.VariantID {
						already = true
						break
					}
				}
				if !already {
					productIDVariantIDs[item.ProductID] = append(variants, item.VariantID)
				}
			}
		}

		// گرفتن دسته‌ی محصول
		p, err := s.pRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return discountservice.EnrichedCart{}, err
		}
		if p.CategoryID != "" {
			categoryIDSet[p.CategoryID] = struct{}{}
		}
	}

	// تبدیل set دسته‌ها به slice
	cids := make([]string, 0, len(categoryIDSet))
	for cid := range categoryIDSet {
		cids = append(cids, cid)
	}

	res := discountservice.EnrichedCart{
		Cart:        cart,
		ProductIDs:  productIDVariantIDs, // map[productID][]variantID
		CategoryIDs: cids,
		Subtotal:    int64(cart.CalculateSubtotal()),
	}

	return res, nil
}

func (s *StoreFrontService) UpdateCartStatus(ctx context.Context, req cartparam.UpdateCartStatusRequest) (cartparam.UpdateCartStatusResponse, error) {
	const op = "cartservice.UpdateCartStatus"
	if err := s.repo.UpdateStatus(ctx, req.CartID, req.Status); err != nil {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.UpdateCartStatusResponse{}, nil
}

func (s *StoreFrontService) AcceptPriceChange(ctx context.Context, req cartparam.AcceptPriceChangeRequest) (cartparam.AcceptPriceChangeResponse, error) {
	const op = "cartservice.AcceptPriceChange"
	variants, err := s.pRepo.GetVariantsByIDs(ctx, []string{req.VariantID})
	if err != nil {
		return cartparam.AcceptPriceChangeResponse{}, richerror.New(op).WithWarpError(err)
	}

	if err := s.repo.UpdateItemPrice(ctx, req.UserID, req.VariantID, int64(variants[req.VariantID].Price)); err != nil {
		return cartparam.AcceptPriceChangeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.AcceptPriceChangeResponse{}, nil
}

func (s *StoreFrontService) DeleteItemFromCart(ctx context.Context, req cartparam.DeleteItemFromCartRequest) (cartparam.DeleteItemFromCartResponse, error) {
	const op = "cartservice.DeleteItemFromCart"

	if err := s.repo.DeleteItemFromCart(ctx, req.UserID, req.VariantID); err != nil {
		return cartparam.DeleteItemFromCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.DeleteItemFromCartResponse{}, nil

}

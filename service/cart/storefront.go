package cartservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StoreFrontService struct {
	repo  CartReadWriter
	pRepo PRepository
	uRepo URepository
}

func NewStorefrontService(repo CartReadWriter,
	pRepo PRepository, uRepo URepository) StoreFrontService {
	return StoreFrontService{
		repo:  repo,
		pRepo: pRepo,
		uRepo: uRepo,
	}
}

func (s StoreFrontService) AddItemToCart(ctx context.Context, req cartparam.AddItemToCartRequest) (cartparam.AddItemToCartResponse, error) {
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

func (s StoreFrontService) RemoveFromCart(ctx context.Context, req cartparam.RemoveFromCartRequest) (cartparam.RemoveFromCartResponse, error) {
	const op = "cartservice.AddItemToCart"
	userID, _ := ctx.Value("user_id").(string)
	if err := s.repo.RemoveFromCart(ctx, userID, req.ProductID, req.VariantID); err != nil {
		return cartparam.RemoveFromCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.RemoveFromCartResponse{}, nil
}

func (s StoreFrontService) GetMyActiveCart(ctx context.Context, req cartparam.MyActiveCartRequest) (cartparam.MyActiveCartResponse, error) {
	const op = "cartservice.GetMyActiveCart"
	cart, err := s.repo.GetUserActiveCart(ctx, req.UserID)
	if err != nil {
		return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	if cart == nil {
		return cartparam.MyActiveCartResponse{}, nil
	}
	ids := []string{req.UserID}
	user, err := s.uRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	var cartItems []cartparam.StoreFrontCartItemDetail
	for _, item := range cart.Items {
		p, err := s.pRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return cartparam.MyActiveCartResponse{}, richerror.New(op).WithWarpError(err)
		}
		cartItems = append(cartItems, cartparam.StoreFrontCartItemDetail{
			ProductID:         p.ID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			Url:               p.GetMainGallery().URL,
		})
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

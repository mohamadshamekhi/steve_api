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
}

func NewStorefrontService(repo CartReadWriter,
	pRepo PRepository) StoreFrontService {
	return StoreFrontService{
		repo:  repo,
		pRepo: pRepo,
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

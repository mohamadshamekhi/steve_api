package cartservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	"golang.org/x/sync/errgroup"
)

type BackofficeService struct {
	repo            CartBackofficeReadWriter
	uRepo           URepository
	pRepo           PRepository
	orderProvider   orderservice.BackofficeService
	productProvider ProductProvider
}

func NewBackofficeService(repo CartBackofficeReadWriter, uRepo URepository, pRepo PRepository, productProvider ProductProvider) BackofficeService {
	return BackofficeService{repo: repo, uRepo: uRepo, pRepo: pRepo, productProvider: productProvider}
}

func (s *BackofficeService) SetOrderProvider(orderProvider orderservice.BackofficeService) {
	s.orderProvider = orderProvider
}

func (s *BackofficeService) GetAll(ctx context.Context, req cartparam.GetAllCartsRequest) (cartparam.GetAllCartsResponse, error) {
	const op = "cartservice.GetAll"

	resp, err := s.repo.GetAll(ctx, req.Filter, req.Pagination)
	if err != nil {
		return cartparam.GetAllCartsResponse{}, richerror.New(op).WithWarpError(err)
	}

	if resp.TotalCount == 0 {
		return cartparam.GetAllCartsResponse{
			Pagination: commonparam.PaginationResponse{
				PageSize:   req.Pagination.GetPageSize(),
				PageNumber: req.Pagination.GetPageNumber(),
				Total:      0,
			},
			Carts: nil,
		}, nil
	}
	userIDs := make([]string, 0, len(resp.Carts))
	for _, cart := range resp.Carts {
		if cart.UserID != "" {
			userIDs = append(userIDs, cart.UserID)
		}
	}

	users, err := s.uRepo.FindManyByIDs(ctx, userIDs)
	if err != nil {
		logger.L().Error("cann't get users , users ids: %v - and errors : %v", userIDs, err)
	}

	var list []cartparam.BackofficeCartList

	for _, cart := range resp.Carts {
		u := entity.User{}
		for _, user := range users {
			if user.ID == cart.UserID {
				u = user
			}
		}
		list = append(list, cartparam.BackofficeCartList{
			ID:              cart.ID,
			Status:          cart.Status,
			UserID:          cart.UserID,
			UserPhoneNumber: u.PhoneNumber,
			UserName:        u.Name,
		})
	}

	return cartparam.GetAllCartsResponse{
		Carts: list,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.PageSize, PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil
}

func (s *BackofficeService) GetSingle(ctx context.Context, req cartparam.GetCartRequest) (cartparam.GetCartResponse, error) {
	const op = "cartservice.GetSingle"
	cart, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	ids := []string{cart.UserID}
	user, err := s.uRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	var cartItems []cartparam.CartItemDetail
	for _, item := range cart.Items {
		p, err := s.pRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return cartparam.GetCartResponse{}, richerror.New(op).WithWarpError(err)
		}
		cartItems = append(cartItems, cartparam.CartItemDetail{
			ProductID:         p.ID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			Url:               p.GetMainGallery().URL,
		})
	}

	c := cartparam.CartDetail{
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
	}
	return cartparam.GetCartResponse{
		Cart: c,
	}, nil
}

func (s *BackofficeService) AddItemToCartBackoffice(ctx context.Context, req cartparam.AddItemToCartBackOfficeRequest) (cartparam.AddItemToCartBackOfficeResponse, error) {
	const op = "cartservice.AddItemToCartBackoffice"
	p, pErr := s.pRepo.FindByID(ctx, req.ProductID)
	if pErr != nil {
		return cartparam.AddItemToCartBackOfficeResponse{}, richerror.New(op).WithWarpError(pErr)
	}
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
	if err := s.repo.AddToCart(ctx, req.UserID, item); err != nil {
		return cartparam.AddItemToCartBackOfficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.AddItemToCartBackOfficeResponse{}, nil
}

func (s *BackofficeService) RemoveFromCartBackoffice(ctx context.Context, req cartparam.RemoveFromCartBackofficeRequest) (cartparam.RemoveFromCartBackofficeResponse, error) {
	const op = "cartservice.RemoveFromCartBackoffice"
	if err := s.repo.RemoveFromCart(ctx, req.UserID, req.ProductID, req.VariantID); err != nil {
		return cartparam.RemoveFromCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return cartparam.RemoveFromCartBackofficeResponse{}, nil
}

func (s *BackofficeService) GetMyActiveCartBackoffice(ctx context.Context, req cartparam.MyActiveCartBackofficeRequest) (cartparam.MyActiveCartBackofficeResponse, error) {
	const op = "cartservice.GetMyActiveCartBackoffice"
	cart, err := s.repo.GetUserActiveCart(ctx, req.UserID)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return cartparam.MyActiveCartBackofficeResponse{}, nil
		}
		return cartparam.MyActiveCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	var cartItems []cartparam.BackofficeCartItemDetail
	for _, item := range cart.Items {
		p, pErr := s.pRepo.FindByID(ctx, item.ProductID)
		if pErr != nil {
			return cartparam.MyActiveCartBackofficeResponse{}, richerror.New(op).WithWarpError(pErr)
		}
		cartItems = append(cartItems, cartparam.BackofficeCartItemDetail{
			ProductID:         p.ID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			ProductName:       item.ProductName,
			VariantAttributes: item.VariantAttributes,
			PriceAtTime:       item.PriceAtTime,
			Url:               p.GetMainGallery().URL,
		})
	}

	c := cartparam.BackofficeCartDetail{
		ID:        cart.ID,
		Status:    string(cart.Status),
		UserID:    cart.UserID,
		Items:     cartItems,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: time.Now(),
		Subtotal:  cart.CalculateSubtotal(),
	}
	return cartparam.MyActiveCartBackofficeResponse{
		Cart: c,
	}, nil
}

func (s *BackofficeService) DeleteItemFromCartBackoffice(ctx context.Context, req cartparam.DeleteItemFromCartBackofficeRequest) (cartparam.DeleteItemFromCartBackofficeResponse, error) {
	const op = "cartservice.DeleteItemFromCart"

	if err := s.repo.DeleteItemFromCart(ctx, req.UserID, req.VariantID); err != nil {
		return cartparam.DeleteItemFromCartBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.DeleteItemFromCartBackofficeResponse{}, nil

}

func (s *BackofficeService) CreateCartBackoffice(ctx context.Context, req cartparam.CreateCartRequest) (cartparam.CreateCartResponse, error) {
	const op = "cartservice.CreateCartBackoffice"
	cart := entity.Cart{
		UserID: req.UserID,
	}
	cart.BeforeCreate()

	err := s.repo.CreateCart(ctx, cart)
	if err != nil {
		return cartparam.CreateCartResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.CreateCartResponse{}, nil
}

func (s *BackofficeService) UpdateStatus(ctx context.Context, req cartparam.UpdateCartStatusRequest) (cartparam.UpdateCartStatusResponse, error) {
	const op = "cartservice.UpdateStatus"

	cart, err := s.repo.FindByID(ctx, req.CartID)
	if err != nil {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !entity.CanTransitionCartStatus(cart.Status, req.Status) {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("cart.invalid_status_transition").
			WithMeta(map[string]any{
				"from": cart.Status,
				"to":   req.Status,
			})
	}
	if uErr := s.repo.UpdateStatus(ctx, req.CartID, req.Status); uErr != nil {
		return cartparam.UpdateCartStatusResponse{}, richerror.New(op).WithWarpError(err)
	}
	return cartparam.UpdateCartStatusResponse{}, nil
}

func (s BackofficeService) CancelExpiredPendingCart(ctx context.Context, req cartparam.CancelExpiredPendingCartRequest) (cartparam.CancelExpiredPendingCartResponse, error) {
	const op = "orderservice.CancelExpiredPendingOrders"

	expiredBefore := time.Now().Add(-req.CanselExpiredPendingCart)

	carts, err := s.repo.GetExpiredPendingCart(ctx, expiredBefore)
	if err != nil {
		return cartparam.CancelExpiredPendingCartResponse{}, richerror.New(op).WithWarpError(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, cart := range carts {
		cartID := cart.ID
		g.Go(func() error {
			var cErr error
			for _, item := range cart.Items {
				_, cErr = s.productProvider.UpdateVariantQuantity(gCtx,
					productparam.UpdateProductVariantQuantityAfterCancelOrderRequest{
						ProductID: item.ProductID,
						VariantID: item.VariantID,
						Quantity:  uint(item.Quantity),
					})
				if cErr != nil {
					logger.L().Error("failed to update variant quantity", "Variant_id", item.VariantID, "err", cErr)
					return cErr
				}
			}

			if cErr = s.repo.UpdateStatus(gCtx, cartID, entity.CartCanceled); err != nil {
				logger.L().Error("failed to update status", "cart_id", cartID, "err", cErr)
				return cErr
			}
			return nil
		})
	}
	_ = g.Wait()

	return cartparam.CancelExpiredPendingCartResponse{}, nil
}

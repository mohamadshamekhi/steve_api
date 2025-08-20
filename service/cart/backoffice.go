package cartservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo  CartBackofficeReadWriter
	uRepo URepository
}

func NewBackofficeService(repo CartBackofficeReadWriter, uRepo URepository) BackofficeService {
	return BackofficeService{repo: repo, uRepo: uRepo}
}

func (s BackofficeService) GetAll(ctx context.Context, req cartparam.GetAllCartsRequest) (cartparam.GetAllCartsResponse, error) {
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

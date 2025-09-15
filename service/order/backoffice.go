package orderservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo         OrderReadWriter
	userProvider UserProvider
}

func NewBackofficeService(repo OrderReadWriter, userProvider UserProvider) BackofficeService {
	return BackofficeService{repo: repo, userProvider: userProvider}
}

func (s BackofficeService) GetAll(ctx context.Context, req orderparam.GetAllOrderRequest) (orderparam.GetAllOrderResponse, error) {
	const op = "orderservice.GetAll"
	resp, err := s.repo.GetAll(ctx, req.Filter, req.Pagination)
	if err != nil {
		return orderparam.GetAllOrderResponse{}, richerror.New(op).WithWarpError(err)
	}
	if resp.TotalCount == 0 {
		return orderparam.GetAllOrderResponse{
			Pagination: commonparam.PaginationResponse{
				PageSize:   req.Pagination.GetPageSize(),
				PageNumber: req.Pagination.GetPageNumber(),
				Total:      0,
			},
			Orders: nil,
		}, nil
	}
	userIDs := make([]string, 0, len(resp.Orders))
	for _, cart := range resp.Orders {
		if cart.UserID != "" {
			userIDs = append(userIDs, cart.UserID)
		}
	}

	users, err := s.userProvider.GetUsersByIDs(ctx, userparam.GetUsersByIDsRequest{IDs: userIDs})
	if err != nil {
		logger.L().Error("cann't get users , users ids: %v - and errors : %v", userIDs, err)
	}

	var list []orderparam.BackofficeOrderList

	for _, order := range resp.Orders {
		u := entity.User{}
		for _, user := range users.Users {
			if user.ID == order.UserID {
				u = entity.User{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					CreatedAt:   user.CreatedAt,
					UpdatedAt:   user.CreatedAt,
					Addresses:   user.Addresses,
				}
			}
		}
		list = append(list, orderparam.BackofficeOrderList{
			ID:              order.ID,
			Status:          order.Status,
			UserID:          order.UserID,
			UserPhoneNumber: u.PhoneNumber,
			UserName:        u.Name,
			SubtotalAmount:  order.SubtotalAmount,
			FinalAmount:     order.FinalAmount,
		})
	}

	return orderparam.GetAllOrderResponse{
		Orders: list,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.PageSize, PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil
}

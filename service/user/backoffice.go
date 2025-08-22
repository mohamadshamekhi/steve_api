package userservice

import (
	"context"

	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type BackofficeService struct {
	repo Repository
}

func NewBackofficeService(repo Repository) BackofficeService {
	return BackofficeService{repo: repo}
}

func (s BackofficeService) GetAllUsers(ctx context.Context, req userparam.GetAllUsersRequest) (userparam.GetAllUsersResponse, error) {
	const op = "userservice.GetAllUsers"

	resp, err := s.repo.GetAllUsers(ctx, req.Filter, req.Pagination)
	if err != nil {
		return userparam.GetAllUsersResponse{}, richerror.New(op).WithWarpError(err)
	}
	var usersInfo []userparam.UserInfo
	for _, u := range resp.Users {
		usersInfo = append(usersInfo, userparam.UserInfo{
			ID:          u.ID,
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
			CreatedAt:   u.CreatedAt,
			Addresses:   u.Addresses,
		})
	}

	return userparam.GetAllUsersResponse{
		Users: usersInfo,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil

}

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

func (s BackofficeService) GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error) {
	const op = "userservice.GetUsersByIDs"
	users, err := s.repo.GetUsersByIDs(ctx, req.IDs)
	if err != nil {
		return userparam.GetUsersByIDsResponse{}, richerror.New(op).WithWarpError(err)
	}
	var userInfos []userparam.UserInfo
	for _, u := range users {
		userInfos = append(userInfos, userparam.UserInfo{
			ID:          u.ID,
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
			CreatedAt:   u.CreatedAt,
			Addresses:   u.Addresses,
		})
	}
	return userparam.GetUsersByIDsResponse{Users: userInfos}, nil
}

func (s BackofficeService) UpsertUser(ctx context.Context, req userparam.UpsertUserRequest) (userparam.UpsertUserResponse, error) {
	const op = "userservice.UpsertUser"

	user, err := s.repo.GetOrCreateUser(ctx, req.PhoneNumber)
	if err != nil {
		return userparam.UpsertUserResponse{}, richerror.New(op).WithWarpError(err)
	}
	nUser := userparam.UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		Addresses:   user.Addresses,
	}
	return userparam.UpsertUserResponse{User: nUser}, nil
}

func (s BackofficeService) UpdateUser(ctx context.Context, req userparam.UpdateUserBackofficeRequest) (userparam.UpdateUserBackofficeResponse, error) {
	const op = "userservice.UpdateUser"
	u, gErr := s.repo.GetUserByID(ctx, req.UserID)
	if gErr != nil {
		return userparam.UpdateUserBackofficeResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	uErr := s.repo.UpdateUser(ctx, u)
	if uErr != nil {
		return userparam.UpdateUserBackofficeResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	return userparam.UpdateUserBackofficeResponse{}, nil
}

func (s BackofficeService) GetUserByID(ctx context.Context, req userparam.GetUsersRequest) (userparam.GetUsersResponse, error) {
	const op = "userservice.GetUserByID"

	user, err := s.repo.GetUserByID(ctx, req.ID)
	if err != nil {
		return userparam.GetUsersResponse{}, richerror.New(op).WithWarpError(err)
	}
	userInfo := userparam.UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		Addresses:   user.Addresses,
	}

	return userparam.GetUsersResponse{User: userInfo}, nil
}

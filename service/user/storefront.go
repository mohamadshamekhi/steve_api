package userservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	userauthparam "github.com/abolfazlnorzad/steve/param/user/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	"github.com/oklog/ulid/v2"
)

type StoreFrontService struct {
	authSvc authservice.Service
	config  Config
	repo    Repository
}

func NewStoreFrontService(config Config, repo Repository, authSvc authservice.Service) StoreFrontService {
	return StoreFrontService{
		config:  config,
		repo:    repo,
		authSvc: authSvc,
	}
}

func (s StoreFrontService) SendOTPCode(ctx context.Context, req userauthparam.UserSendLoginCodeRequest) (userauthparam.UserSendLoginCodeResponse, error) {
	const op = "userservice.SendOTPCode"

	err := s.authSvc.GenerateOTPCode(ctx, req.PhoneNumber, s.config.Prefix)
	if err != nil {
		return userauthparam.UserSendLoginCodeResponse{}, richerror.New(op)
	}

	// return resp.
	return userauthparam.UserSendLoginCodeResponse{}, nil
}

func (s StoreFrontService) VerifyOTPCode(ctx context.Context, req userauthparam.UserVerifyLoginCodeRequest) (userauthparam.UserVerifyLoginCodeResponse, error) {
	const op = "userservice.VerifyOTPCode"

	resp, err := s.authSvc.VerifyOTPCode(ctx, req.PhoneNumber, req.Code, s.config.Prefix)
	if err != nil {
		return userauthparam.UserVerifyLoginCodeResponse{}, err
	}
	if resp {
		user, iErr := s.repo.GetOrCreateUser(ctx, req.PhoneNumber)
		if iErr != nil {
			return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		token, gErr := s.authSvc.GenerateAccessToken(user.ID, s.config.Prefix, []string{}, []string{})
		if gErr != nil {
			return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(gErr)
		}
		return userauthparam.UserVerifyLoginCodeResponse{
			Token: token,
		}, nil
	}

	return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).
		WithMessageKey("validation.incorrect_opt_code")
}

func (s StoreFrontService) UpdateProfile(ctx context.Context, req userparam.StoreFrontUpdateProfileRequest) (userparam.StoreFrontUpdateProfileResponse, error) {
	const op = "userservice.UpdateProfile"

	u, gErr := s.repo.GetUserByID(ctx, req.UserID)
	if gErr != nil {
		return userparam.StoreFrontUpdateProfileResponse{}, richerror.New(op).WithWarpError(gErr)
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		u.PhoneNumber = *req.PhoneNumber
	}
	uErr := s.repo.UpdateUser(ctx, u)
	if uErr != nil {
		return userparam.StoreFrontUpdateProfileResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	return userparam.StoreFrontUpdateProfileResponse{}, nil
}

func (s StoreFrontService) AddAddress(ctx context.Context, req userparam.StoreFrontAddAddressRequest) (userparam.StoreFrontAddAddressResponse, error) {
	const op = "userservice.AddAddress"

	if err := s.repo.AddAddress(ctx, req.UserID, entity.Address{
		ID:       ulid.Make().String(),
		Title:    req.Title,
		Province: req.Province,
		City:     req.City,
		Street:   req.Street,
		Postcode: req.Postcode,
	}); err != nil {
		return userparam.StoreFrontAddAddressResponse{}, richerror.New(op).WithWarpError(err)
	}
	return userparam.StoreFrontAddAddressResponse{}, nil
}

func (s StoreFrontService) UpdateAddress(ctx context.Context, req userparam.StoreFrontUpdateAddressRequest) (userparam.StoreFrontUpdateAddressResponse, error) {
	const op = "userservice.UpdateAddress"
	ad, err := s.repo.FindAddress(ctx, req.UserID, req.ID)
	if err != nil {
		return userparam.StoreFrontUpdateAddressResponse{}, richerror.New(op).WithWarpError(err)
	}
	if req.City != nil {
		ad.City = *req.City
	}
	if req.Title != nil {
		ad.Title = *req.Title
	}
	if req.Province != nil {
		ad.Province = *req.Province
	}
	if req.Street != nil {
		ad.Street = *req.Street
	}
	if req.Postcode != nil {
		ad.Postcode = *req.Postcode
	}

	if err = s.repo.UpdateAddress(ctx, req.UserID, ad); err != nil {
		return userparam.StoreFrontUpdateAddressResponse{}, richerror.New(op).WithWarpError(err)
	}
	return userparam.StoreFrontUpdateAddressResponse{}, nil
}

func (s StoreFrontService) FindAddressByID(ctx context.Context, req userparam.FindAddressByIDRequest) (userparam.FindAddressByIDResponse, error) {
	const op = "userservice.FindAddressByID"
	address, err := s.repo.FindAddress(ctx, req.UserID, req.AddressID)
	if err != nil {
		return userparam.FindAddressByIDResponse{}, richerror.New(op).WithWarpError(err)
	}
	return userparam.FindAddressByIDResponse{Address: address}, nil
}

package userservice

import (
	"context"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	userauthparam "github.com/abolfazlnorzad/steve/param/user/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/phonenumber"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	"github.com/oklog/ulid/v2"
)

type StoreFrontService struct {
	authSvc       authservice.Service
	config        Config
	repo          Repository
	cacheProvider CacheProvider
	referral      ReferralAttributor
}
type CacheProvider interface {
	AddToBlocklist(ctx context.Context, kind string, JWTID string, ExpireAt time.Time) error
}

func NewStoreFrontService(config Config, repo Repository, authSvc authservice.Service, cacheProvider CacheProvider, referral ReferralAttributor) StoreFrontService {
	return StoreFrontService{
		config:        config,
		repo:          repo,
		authSvc:       authSvc,
		cacheProvider: cacheProvider,
		referral:      referral,
	}
}

func (s StoreFrontService) SendOTPCode(ctx context.Context, req userauthparam.UserSendLoginCodeRequest) (userauthparam.UserSendLoginCodeResponse, error) {
	const op = "userservice.SendOTPCode"

	req.PhoneNumber = phonenumber.Normalize(req.PhoneNumber)

	err := s.authSvc.GenerateOTPCode(ctx, req.PhoneNumber, s.config.Prefix)
	if err != nil {
		return userauthparam.UserSendLoginCodeResponse{}, richerror.New(op)
	}

	// Tell the frontend whether the next step is a sign-up (show the ref field)
	// or a login. Non-fatal: on error assume existing user.
	exists, eErr := s.repo.CheckUserIsExistByPhoneNumber(ctx, req.PhoneNumber)
	if eErr != nil {
		logger.L().Error("userservice: new-user check failed", "err", eErr)
		exists = true
	}

	return userauthparam.UserSendLoginCodeResponse{IsNewUser: !exists}, nil
}

func (s StoreFrontService) VerifyOTPCode(ctx context.Context, req userauthparam.UserVerifyLoginCodeRequest) (userauthparam.UserVerifyLoginCodeResponse, error) {
	const op = "userservice.VerifyOTPCode"

	req.PhoneNumber = phonenumber.Normalize(req.PhoneNumber)

	resp, err := s.authSvc.VerifyOTPCode(ctx, req.PhoneNumber, req.Code, s.config.Prefix)
	if err != nil {
		return userauthparam.UserVerifyLoginCodeResponse{}, err
	}
	if resp {
		// Detect a brand-new signup before the upsert so we can credit a referral.
		existed, eErr := s.repo.CheckUserIsExistByPhoneNumber(ctx, req.PhoneNumber)
		if eErr != nil {
			logger.L().Error("userservice: referral new-user check failed", "err", eErr)
			existed = true // fail safe: skip attribution, never block login
		}

		var birthDate *time.Time
		if !existed {
			// Format already enforced by the validator for new users.
			if t, pErr := time.Parse(constant.BirthDateLayout, req.BirthDate); pErr == nil {
				birthDate = &t
			}
		}
		user, iErr := s.repo.GetOrCreateUserWithProfile(ctx, req.PhoneNumber, strings.TrimSpace(req.Name), birthDate)
		if iErr != nil {
			return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		token, gErr := s.authSvc.GenerateAccessToken(user.ID, s.config.Prefix, []string{}, []string{})
		if gErr != nil {
			return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(gErr)
		}

		if !existed && s.referral != nil {
			s.attributeReferralAsync(user.ID, req.PhoneNumber, req.Ref)
		}

		return userauthparam.UserVerifyLoginCodeResponse{
			Token:     token,
			IsNewUser: !existed,
		}, nil
	}

	return userauthparam.UserVerifyLoginCodeResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).
		WithMessageKey("validation.incorrect_opt_code")
}

// attributeReferralAsync runs the referral signup hook off the request
// goroutine with its own context so a slow SMS or coupon write never adds
// latency to — or fails — the login response.
func (s StoreFrontService) attributeReferralAsync(userID, phoneNumber, refCode string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.referral.AttributeSignup(ctx, userID, phoneNumber, refCode); err != nil {
			logger.L().Error("userservice: referral attribution failed", "err", err, "user_id", userID)
		}
	}()
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

func (s StoreFrontService) Logout(ctx context.Context, req userparam.LogoutRequest) (userparam.LogoutResponse, error) {
	const op = "userservice.Logout"
	if err := s.cacheProvider.AddToBlocklist(ctx, req.Kind, req.JWTID, req.ExpireAt); err != nil {
		return userparam.LogoutResponse{}, richerror.New(op).WithWarpError(err)
	}

	return userparam.LogoutResponse{}, nil
}

func (s StoreFrontService) Me(ctx context.Context, request userparam.MeRequest) (userparam.MeResponse, error) {
	const OP = "userservice.Me"
	u, err := s.repo.GetUsersByIDs(ctx, []string{request.UserID})
	if err != nil {
		return userparam.MeResponse{}, richerror.New(OP).WithWarpError(err)
	}
	if len(u) == 0 {
		return userparam.MeResponse{}, richerror.New(OP).WithKind(richerror.KindUnauthorized).WithMessageKey("validation.unauthorized")
	}
	userInfo := userparam.UserInfo{
		ID:          u[0].ID,
		Name:        u[0].Name,
		PhoneNumber: u[0].PhoneNumber,
		BirthDate:   u[0].BirthDate,
		CreatedAt:   u[0].CreatedAt,
		Addresses:   u[0].Addresses,
	}

	return userparam.MeResponse{
		User: userInfo,
	}, nil

}

func (s StoreFrontService) GetUserByID(ctx context.Context, req userparam.GetUsersRequest) (userparam.GetUsersResponse, error) {
	const OP = "userservice.GetUserByID"

	user, err := s.repo.GetUserByID(ctx, req.ID)
	if err != nil {
		return userparam.GetUsersResponse{}, richerror.New(OP).WithWarpError(err)
	}
	userInfo := userparam.UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		BirthDate:   user.BirthDate,
		CreatedAt:   user.CreatedAt,
		Addresses:   user.Addresses,
	}

	return userparam.GetUsersResponse{User: userInfo}, nil
}

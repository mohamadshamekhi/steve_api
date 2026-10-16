package userservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type Config struct {
	Prefix string `koanf:"prefix"`
}

// ReferralAttributor is the signup-time hook that credits an invite when a
// brand-new user finishes OTP registration. Best-effort: the storefront service
// runs it off the request goroutine and never fails a login on its error.
// Satisfied by referralservice.StorefrontService.
type ReferralAttributor interface {
	AttributeSignup(ctx context.Context, newUserID, phoneNumber, refCode string) error
}

type ListResult struct {
	Users      []entity.User
	TotalCount int64
}

type Repository interface {
	GetUserByPhoneNumber(ctx context.Context, phone string) (entity.User, error)
	GetUserByID(ctx context.Context, phone string) (entity.User, error)
	AddNewUser(ctx context.Context, u entity.User) error
	GetOrCreateUser(ctx context.Context, phonenumber string) (entity.User, error)
	CheckUserIsExistByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	UpdateUser(ctx context.Context, u entity.User) error
	AddAddress(ctx context.Context, userID string, address entity.Address) error
	UpdateAddress(ctx context.Context, userID string, address entity.Address) error
	FindAddress(ctx context.Context, userID string, addressID string) (entity.Address, error)
	GetAllUsers(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error)
}

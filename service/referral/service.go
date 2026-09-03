package referralservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// Repository is the referral persistence port (three collections behind one
// *mongodbreferral.DB).
type Repository interface {
	FindAccountByUserID(ctx context.Context, userID string) (entity.ReferralAccount, error)
	FindAccountByCode(ctx context.Context, code string) (entity.ReferralAccount, error)
	CreateAccount(ctx context.Context, a entity.ReferralAccount) error
	IncrementAccountCounters(ctx context.Context, userID string, joinedDelta, rewardedDelta int) (entity.ReferralAccount, error)

	FindReferralByReferredUserID(ctx context.Context, referredUserID string) (entity.Referral, error)
	FindPendingInviteByPhone(ctx context.Context, phone string) (entity.Referral, error)
	FindPendingInviteByReferrerAndPhone(ctx context.Context, referrerUserID, phone string) (entity.Referral, error)
	ListUnrewardedJoined(ctx context.Context, referrerUserID string, limit int) ([]entity.Referral, error)
	ListReferralsByReferrer(ctx context.Context, referrerUserID string, p commonparam.PaginationRequest) ([]entity.Referral, int64, error)
	ListRewardsByReferrer(ctx context.Context, referrerUserID string) ([]entity.ReferralReward, error)

	CreateReferral(ctx context.Context, r entity.Referral) error
	CreateReward(ctx context.Context, rw entity.ReferralReward) error
	MarkReferralJoined(ctx context.Context, referralID, referredUserID string) error
	MarkReferralsRewarded(ctx context.Context, referralIDs []string, rewardBatchID string) error
}

// DiscountProvider is the slice of the discount repo the referral domain needs
// to mint reward / welcome coupons. Satisfied by *mongodbdiscount.DB.
type DiscountProvider interface {
	Create(ctx context.Context, discount entity.Discount) error
	CheckCodeIsUnique(ctx context.Context, code string) (bool, error)
}

// UserProvider resolves invitee names/phones for the "who joined" list.
// Satisfied by *mongodbuser.DB.
type UserProvider interface {
	FindManyByIDs(ctx context.Context, ids []string) ([]entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

// SMSSender sends a free-text SMS (invite links + reward notifications).
type SMSSender interface {
	SendSimple(ctx context.Context, to string, message string) error
}

// StorefrontService is the single service for this domain: it powers the
// storefront "invite friends" screen and the signup-time attribution hook.
type StorefrontService struct {
	cfg          Config
	repo         Repository
	discountRepo DiscountProvider
	userRepo     UserProvider
	sms          SMSSender
}

func NewStorefrontService(cfg Config, repo Repository, discountRepo DiscountProvider, userRepo UserProvider, sms SMSSender) StorefrontService {
	return StorefrontService{
		cfg:          cfg,
		repo:         repo,
		discountRepo: discountRepo,
		userRepo:     userRepo,
		sms:          sms,
	}
}

package birthdayservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// Repository is the birthday persistence port. Satisfied by *mongodbbirthday.DB.
type Repository interface {
	GetSettings(ctx context.Context) (entity.BirthdaySettings, error)
	UpsertSettings(ctx context.Context, s entity.BirthdaySettings) error

	HasRewardSince(ctx context.Context, userID string, since time.Time) (bool, error)
	ClaimReward(ctx context.Context, r entity.BirthdayReward) (bool, error)
	ReleaseReward(ctx context.Context, id string) error
	SetRewardDiscount(ctx context.Context, id, discountID, discountCode string) error
	MarkRewardSMSSent(ctx context.Context, id string) error
	ListRewards(ctx context.Context, p commonparam.PaginationRequest) ([]entity.BirthdayReward, int64, error)
}

// UserProvider finds users by birthday. Satisfied by *mongodbuser.DB.
type UserProvider interface {
	FindUsersByBirthMonthDay(ctx context.Context, days []entity.MonthDay) ([]entity.User, error)
}

// DiscountProvider mints the birthday coupons. Satisfied by *mongodbdiscount.DB.
type DiscountProvider interface {
	Create(ctx context.Context, discount entity.Discount) error
	CheckCodeIsUnique(ctx context.Context, code string) (bool, error)
}

// SMSSender sends a template ("pattern") SMS. Satisfied by *mellipayamak.Adapter.
type SMSSender interface {
	SendPatternSMS(ctx context.Context, to string, text []string, patternKey string) error
}

// Service runs the daily birthday job and backs the backoffice settings screen.
type Service struct {
	cfg          Config
	repo         Repository
	userRepo     UserProvider
	discountRepo DiscountProvider
	sms          SMSSender
}

func New(cfg Config, repo Repository, userRepo UserProvider, discountRepo DiscountProvider, sms SMSSender) Service {
	return Service{
		cfg:          cfg,
		repo:         repo,
		userRepo:     userRepo,
		discountRepo: discountRepo,
		sms:          sms,
	}
}

// effective returns the admin-saved settings, falling back to config defaults
// when none were saved (or the read fails). custom reports which one it is.
func (s Service) effective(ctx context.Context) (st entity.BirthdaySettings, custom bool, err error) {
	st, err = s.repo.GetSettings(ctx)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			return s.cfg.toSettings(), false, nil
		}
		return s.cfg.toSettings(), false, err
	}
	return st, true, nil
}

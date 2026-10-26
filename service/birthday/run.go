package birthdayservice

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// tehranTZ is fixed at +03:30; Iran has no DST since 2022.
var tehranTZ = time.FixedZone("IRST", 3*3600+30*60)

// Run grants birthday codes to every user whose birthday falls between today
// and today+DaysBefore (Tehran time) and who has no code for that birthday yet.
// Covering the whole window, not just the exact day, means a missed run or a
// late signup still gets its code. Safe to run repeatedly.
func (s Service) Run(ctx context.Context) (birthdayparam.RunResponse, error) {
	const op = "birthdayservice.Run"

	st, _, err := s.effective(ctx)
	if err != nil {
		// Don't grant with defaults the admin may have overridden.
		return birthdayparam.RunResponse{}, richerror.New(op).WithWarpError(err)
	}
	resp := birthdayparam.RunResponse{Enabled: st.Enabled}
	if !st.Enabled {
		return resp, nil
	}

	now := time.Now().In(tehranTZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tehranTZ)

	for offset := 0; offset <= daysBefore(st); offset++ {
		day := today.AddDate(0, 0, offset)
		users, fErr := s.userRepo.FindUsersByBirthMonthDay(ctx, matchingDays(day))
		if fErr != nil {
			return resp, richerror.New(op).WithWarpError(fErr)
		}
		for _, u := range users {
			resp.Matched++
			granted, gErr := s.grant(ctx, st, u, day)
			switch {
			case gErr != nil:
				resp.Failed++
				logger.L().Error("birthday: failed to grant code", "err", gErr, "user_id", u.ID)
			case granted:
				resp.Granted++
			default:
				resp.Skipped++
			}
		}
	}
	return resp, nil
}

// matchingDays is the birth month/day that celebrates on day. In a non-leap
// year, people born on Feb 29 celebrate on Feb 28.
func matchingDays(day time.Time) []entity.MonthDay {
	days := []entity.MonthDay{{Month: day.Month(), Day: day.Day()}}
	if day.Month() == time.February && day.Day() == 28 && !isLeap(day.Year()) {
		days = append(days, entity.MonthDay{Month: time.February, Day: 29})
	}
	return days
}

func isLeap(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// grant claims (user, year), mints the coupon and sends the SMS. granted is
// false when this birthday was already rewarded.
func (s Service) grant(ctx context.Context, st entity.BirthdaySettings, u entity.User, birthday time.Time) (bool, error) {
	const op = "birthdayservice.grant"

	r := entity.BirthdayReward{
		UserID:      u.ID,
		UserName:    u.Name,
		PhoneNumber: u.PhoneNumber,
		Year:        birthday.Year(),
		BirthdayOn:  birthday,
	}
	r.BeforeCreate()

	ok, err := s.repo.ClaimReward(ctx, r)
	if err != nil {
		return false, richerror.New(op).WithWarpError(err)
	}
	if !ok {
		return false, nil
	}

	d, mErr := s.mintCoupon(ctx, st, u.ID)
	if mErr != nil {
		if rErr := s.repo.ReleaseReward(ctx, r.ID); rErr != nil {
			logger.L().Error("birthday: failed to release claim", "err", rErr, "user_id", u.ID)
		}
		return false, richerror.New(op).WithWarpError(mErr)
	}
	if uErr := s.repo.SetRewardDiscount(ctx, r.ID, d.ID, d.Code); uErr != nil {
		// Coupon exists and the claim blocks a re-grant; only the audit link is missing.
		logger.L().Error("birthday: failed to store discount on reward", "err", uErr, "user_id", u.ID, "code", d.Code)
	}

	// {0}=name {1}=code {2}=discount label {3}=valid days
	text := []string{displayName(u.Name), d.Code, discountLabel(st), strconv.Itoa(validDays(st))}
	if sErr := s.sms.SendPatternSMS(ctx, u.PhoneNumber, text, mellipayamak.BirthdayDiscountKey); sErr != nil {
		logger.L().Error("birthday: failed to send sms", "err", sErr, "user_id", u.ID)
		return true, nil
	}
	if mErr := s.repo.MarkRewardSMSSent(ctx, r.ID); mErr != nil {
		logger.L().Error("birthday: failed to mark sms sent", "err", mErr, "user_id", u.ID)
	}
	return true, nil
}

func (s Service) mintCoupon(ctx context.Context, st entity.BirthdaySettings, userID string) (entity.Discount, error) {
	const op = "birthdayservice.mintCoupon"

	dtype := entity.DiscountType(st.DiscountType)
	if dtype != entity.PercentageDiscount && dtype != entity.FixedAmountDiscount {
		dtype = entity.FixedAmountDiscount
	}

	var code string
	for attempt := 0; attempt < 6; attempt++ {
		c, gErr := generateCouponCode()
		if gErr != nil {
			return entity.Discount{}, richerror.New(op).WithWarpError(gErr)
		}
		unique, uErr := s.discountRepo.CheckCodeIsUnique(ctx, c)
		if uErr != nil {
			return entity.Discount{}, richerror.New(op).WithWarpError(uErr)
		}
		if unique {
			code = c
			break
		}
	}
	if code == "" {
		return entity.Discount{}, richerror.New(op).
			WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}

	d := entity.Discount{
		Code:                  code,
		Type:                  dtype,
		Value:                 st.DiscountValue,
		MinimumPurchaseAmount: st.MinimumPurchaseAmount,
		StartsAt:              time.Now().Add(-time.Minute),
		ExpiresAt:             time.Now().AddDate(0, 0, validDays(st)),
		TotalUsageLimit:       1,
		UsageLimitPerUser:     1,
		ApplicableTo:          entity.ApplicableToUsers,
		ApplicableUserIDs:     []string{userID},
		IsActive:              true,
	}
	d.BeforeCreate()
	if err := s.discountRepo.Create(ctx, d); err != nil {
		return entity.Discount{}, richerror.New(op).WithWarpError(err)
	}
	return d, nil
}

const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func generateCouponCode() (string, error) {
	const n = 6
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = codeAlphabet[idx.Int64()]
	}
	return "BDAY" + string(b), nil
}

func displayName(name string) string {
	if n := strings.TrimSpace(name); n != "" {
		return n
	}
	return "مشتری"
}

// discountLabel renders the discount for the SMS, e.g. "500,000 تومان" or "10 درصد".
func discountLabel(st entity.BirthdaySettings) string {
	if entity.DiscountType(st.DiscountType) == entity.PercentageDiscount {
		return strconv.FormatUint(uint64(st.DiscountValue), 10) + " درصد"
	}
	return groupThousands(uint64(st.DiscountValue)) + " تومان"
}

func groupThousands(v uint64) string {
	s := strconv.FormatUint(v, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

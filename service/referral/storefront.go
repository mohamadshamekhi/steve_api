package referralservice

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/phonenumber"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// codeAlphabet excludes visually ambiguous characters (0/O, 1/I/L).
const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func (s StorefrontService) programEnabled(op string) error {
	if !s.cfg.Enabled {
		return richerror.New(op).WithKind(richerror.KindForbidden).
			WithMessageKey("referral.program_disabled")
	}
	return nil
}

// effective returns the runtime config: config/env defaults with any
// admin-saved settings document overlaid on top.
func (s StorefrontService) effective(ctx context.Context) Config {
	st, err := s.repo.GetSettings(ctx)
	if err != nil {
		return s.cfg
	}
	return s.cfg.applySettings(st)
}

// GetMyLink returns the caller's permanent invite code + link, creating the
// referral account lazily on first call.
func (s StorefrontService) GetMyLink(ctx context.Context, userID string) (referralparam.GetMyLinkResponse, error) {
	const op = "referralservice.GetMyLink"
	if err := s.programEnabled(op); err != nil {
		return referralparam.GetMyLinkResponse{}, err
	}

	acc, err := s.getOrCreateAccount(ctx, userID)
	if err != nil {
		return referralparam.GetMyLinkResponse{}, richerror.New(op).WithWarpError(err)
	}

	rewards, err := s.repo.ListRewardsByReferrer(ctx, userID)
	if err != nil {
		return referralparam.GetMyLinkResponse{}, richerror.New(op).WithWarpError(err)
	}

	return s.buildLinkResponse(s.effective(ctx), acc, len(rewards)), nil
}

func (s StorefrontService) buildLinkResponse(eff Config, acc entity.ReferralAccount, rewardsEarned int) referralparam.GetMyLinkResponse {
	threshold := eff.rewardThreshold()
	nextIn := threshold - acc.PendingProgress()
	if nextIn <= 0 {
		nextIn = threshold
	}
	return referralparam.GetMyLinkResponse{
		Code:            acc.Code,
		URL:             s.buildURL(acc.Code),
		TotalJoined:     acc.TotalJoined,
		PendingProgress: acc.PendingProgress(),
		RewardThreshold: threshold,
		NextRewardIn:    nextIn,
		RewardsEarned:   rewardsEarned,
	}
}

func (s StorefrontService) buildURL(code string) string {
	base := strings.TrimRight(s.cfg.InviteBaseURL, "/")
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%sref=%s", base, sep, code)
}

// InviteByPhone records an outstanding invite and texts the link to the number.
func (s StorefrontService) InviteByPhone(ctx context.Context, userID string, req referralparam.InviteByPhoneRequest) (referralparam.InviteByPhoneResponse, error) {
	const op = "referralservice.InviteByPhone"
	if err := s.programEnabled(op); err != nil {
		return referralparam.InviteByPhoneResponse{}, err
	}

	phone := normalizePhone(req.PhoneNumber)

	caller, err := s.userRepo.GetUserByID(ctx, userID)
	if err == nil && normalizePhone(caller.PhoneNumber) == phone {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("referral.cannot_invite_self")
	}

	// Can't invite someone who already has an account.
	if exists, uErr := s.userRepo.CheckUserIsExistByPhoneNumber(ctx, phone); uErr != nil {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).WithWarpError(uErr)
	} else if exists {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("referral.already_registered")
	}

	// Can't invite a number that anyone (this referrer or another) already invited.
	if invited, iErr := s.repo.ExistsInviteByPhone(ctx, phone); iErr != nil {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).WithWarpError(iErr)
	} else if invited {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("referral.already_invited")
	}

	acc, err := s.getOrCreateAccount(ctx, userID)
	if err != nil {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).WithWarpError(err)
	}

	url := s.buildURL(acc.Code)

	// Send the SMS *before* recording the invite: if the provider rejects it we
	// don't want a phantom "invited" row that permanently blocks this number from
	// ever being invited again (by this referrer or anyone else). On failure the
	// caller still gets the link to share manually and can retry later.
	// Template: "... کد دعوت {0} ..." + site name.
	if sErr := s.sms.SendPatternSMS(ctx, phone, []string{acc.Code}, mellipayamak.ReferralInviteKey); sErr != nil {
		logger.L().Error("referral: failed to send invite sms", "err", sErr, "phone", phone)
		return referralparam.InviteByPhoneResponse{URL: url, SMSSent: false}, nil
	}

	r := entity.Referral{
		ReferrerUserID: userID,
		Code:           acc.Code,
		Source:         entity.ReferralSourceSMS,
		Status:         entity.ReferralStatusInvited,
		InvitedPhone:   phone,
	}
	r.BeforeCreate()
	if cErr := s.repo.CreateReferral(ctx, r); cErr != nil {
		return referralparam.InviteByPhoneResponse{}, richerror.New(op).WithWarpError(cErr)
	}

	return referralparam.InviteByPhoneResponse{ReferralID: r.ID, URL: url, SMSSent: true}, nil
}

// ListMyReferrals returns the paginated "who joined via my link" list.
func (s StorefrontService) ListMyReferrals(ctx context.Context, userID string, req referralparam.ListMyReferralsRequest) (referralparam.ListMyReferralsResponse, error) {
	const op = "referralservice.ListMyReferrals"
	if err := s.programEnabled(op); err != nil {
		return referralparam.ListMyReferralsResponse{}, err
	}

	acc, err := s.getOrCreateAccount(ctx, userID)
	if err != nil {
		return referralparam.ListMyReferralsResponse{}, richerror.New(op).WithWarpError(err)
	}

	rows, total, err := s.repo.ListReferralsByReferrer(ctx, userID, req.Pagination)
	if err != nil {
		return referralparam.ListMyReferralsResponse{}, richerror.New(op).WithWarpError(err)
	}

	// Resolve names for joined invitees.
	joinedIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.ReferredUserID != "" {
			joinedIDs = append(joinedIDs, r.ReferredUserID)
		}
	}
	usersByID := make(map[string]entity.User)
	if len(joinedIDs) > 0 {
		users, uErr := s.userRepo.FindManyByIDs(ctx, joinedIDs)
		if uErr != nil {
			return referralparam.ListMyReferralsResponse{}, richerror.New(op).WithWarpError(uErr)
		}
		for _, u := range users {
			usersByID[u.ID] = u
		}
	}

	items := make([]referralparam.ReferralItem, 0, len(rows))
	for _, r := range rows {
		item := referralparam.ReferralItem{
			ID:        r.ID,
			Source:    string(r.Source),
			Status:    string(r.Status),
			Rewarded:  r.RewardBatchID != "",
			InvitedAt: r.InvitedAt,
			JoinedAt:  r.JoinedAt,
		}
		if u, ok := usersByID[r.ReferredUserID]; ok {
			item.Name = u.Name
			item.PhoneMasked = maskPhone(u.PhoneNumber)
		} else if r.InvitedPhone != "" {
			item.PhoneMasked = maskPhone(r.InvitedPhone)
		}
		items = append(items, item)
	}

	threshold := s.effective(ctx).rewardThreshold()
	nextIn := threshold - acc.PendingProgress()
	if nextIn <= 0 {
		nextIn = threshold
	}

	return referralparam.ListMyReferralsResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(total),
		},
		Items:           items,
		TotalJoined:     acc.TotalJoined,
		PendingProgress: acc.PendingProgress(),
		RewardThreshold: threshold,
		NextRewardIn:    nextIn,
	}, nil
}

// ListMyRewards returns the discount codes the caller has earned.
func (s StorefrontService) ListMyRewards(ctx context.Context, userID string) (referralparam.ListMyRewardsResponse, error) {
	const op = "referralservice.ListMyRewards"
	if err := s.programEnabled(op); err != nil {
		return referralparam.ListMyRewardsResponse{}, err
	}
	rewards, err := s.repo.ListRewardsByReferrer(ctx, userID)
	if err != nil {
		return referralparam.ListMyRewardsResponse{}, richerror.New(op).WithWarpError(err)
	}
	items := make([]referralparam.RewardItem, 0, len(rewards))
	for _, rw := range rewards {
		items = append(items, referralparam.RewardItem{
			ID:           rw.ID,
			DiscountCode: rw.DiscountCode,
			CreatedAt:    rw.CreatedAt,
		})
	}
	return referralparam.ListMyRewardsResponse{Items: items}, nil
}

// getOrCreateAccount returns the caller's referral account, minting one with a
// unique code on first use.
func (s StorefrontService) getOrCreateAccount(ctx context.Context, userID string) (entity.ReferralAccount, error) {
	const op = "referralservice.getOrCreateAccount"

	acc, err := s.repo.FindAccountByUserID(ctx, userID)
	if err == nil {
		return acc, nil
	}
	if richerror.GetKind(err) != richerror.KindNotFound {
		return entity.ReferralAccount{}, richerror.New(op).WithWarpError(err)
	}

	for attempt := 0; attempt < 5; attempt++ {
		code, gErr := generateCode()
		if gErr != nil {
			return entity.ReferralAccount{}, richerror.New(op).WithWarpError(gErr)
		}
		if _, fErr := s.repo.FindAccountByCode(ctx, code); fErr == nil {
			continue // collision, retry
		}
		newAcc := entity.ReferralAccount{UserID: userID, Code: code}
		newAcc.BeforeCreate()
		if cErr := s.repo.CreateAccount(ctx, newAcc); cErr != nil {
			// Another request created it concurrently — read it back.
			if existing, rErr := s.repo.FindAccountByUserID(ctx, userID); rErr == nil {
				return existing, nil
			}
			continue
		}
		return newAcc, nil
	}
	return entity.ReferralAccount{}, richerror.New(op).
		WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
}

func generateCode() (string, error) {
	const n = 5
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = codeAlphabet[idx.Int64()]
	}
	return "STV-" + string(b), nil
}

func normalizePhone(p string) string {
	return phonenumber.Normalize(p)
}

func maskPhone(p string) string {
	p = strings.TrimSpace(p)
	if len(p) < 7 {
		return "***"
	}
	return p[:4] + "***" + p[len(p)-2:]
}

// mintCoupon builds and persists a single-use, user-scoped discount and returns
// its code.
func (s StorefrontService) mintCoupon(ctx context.Context, userID, discType string, value uint, minPurchase uint, validDays int) (entity.Discount, error) {
	const op = "referralservice.mintCoupon"

	dtype := entity.DiscountType(discType)
	if dtype != entity.PercentageDiscount && dtype != entity.FixedAmountDiscount {
		dtype = entity.PercentageDiscount
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
		Value:                 value,
		MinimumPurchaseAmount: minPurchase,
		StartsAt:              time.Now().Add(-time.Minute),
		ExpiresAt:             time.Now().AddDate(0, 0, validDays),
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

func generateCouponCode() (string, error) {
	const n = 8
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = codeAlphabet[idx.Int64()]
	}
	return "REF" + string(b), nil
}

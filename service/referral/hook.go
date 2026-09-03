package referralservice

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// AttributeSignup is the signup-time hook: called once, right after a brand-new
// storefront user finishes OTP registration. It is best-effort — every failure
// is logged and swallowed so it can never block a login. Callers should run it
// off the request goroutine.
//
//   - refCode:      the ?ref=<code> the invitee arrived with (may be empty)
//   - phoneNumber:  the new user's phone, used to match an outstanding SMS invite
func (s StorefrontService) AttributeSignup(ctx context.Context, newUserID, phoneNumber, refCode string) error {
	const op = "referralservice.AttributeSignup"

	if !s.cfg.Enabled {
		return nil
	}

	phone := normalizePhone(phoneNumber)

	// 1. Resolve which invite code brought this user in.
	code := normalizeCode(refCode)
	var promoteReferralID string
	if code == "" {
		if inv, err := s.repo.FindPendingInviteByPhone(ctx, phone); err == nil {
			code = inv.Code
			promoteReferralID = inv.ID
		}
	}
	if code == "" {
		return nil // organic signup, nothing to attribute
	}

	// 2. Resolve the referrer.
	acc, err := s.repo.FindAccountByCode(ctx, code)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			logger.L().Info("referral: signup with unknown ref code", "code", code, "user_id", newUserID)
			return nil
		}
		return s.logHookErr(op, err, newUserID)
	}

	// 3. Guard: no self-referral.
	if acc.UserID == newUserID {
		return nil
	}

	// 4. Guard: this user was already attributed to someone.
	if _, fErr := s.repo.FindReferralByReferredUserID(ctx, newUserID); fErr == nil {
		return nil
	} else if richerror.GetKind(fErr) != richerror.KindNotFound {
		return s.logHookErr(op, fErr, newUserID)
	}

	// 5. Promote a matching SMS invite, or record a fresh link join.
	if promoteReferralID == "" {
		if inv, iErr := s.repo.FindPendingInviteByPhone(ctx, phone); iErr == nil && inv.Code == code {
			promoteReferralID = inv.ID
		}
	}
	if promoteReferralID != "" {
		if mErr := s.repo.MarkReferralJoined(ctx, promoteReferralID, newUserID); mErr != nil {
			return s.logHookErr(op, mErr, newUserID)
		}
	} else {
		now := time.Now()
		r := entity.Referral{
			ReferrerUserID: acc.UserID,
			Code:           code,
			Source:         entity.ReferralSourceLink,
			Status:         entity.ReferralStatusJoined,
			ReferredUserID: newUserID,
			InvitedAt:      now,
			JoinedAt:       &now,
		}
		r.BeforeCreate()
		if cErr := s.repo.CreateReferral(ctx, r); cErr != nil {
			return s.logHookErr(op, cErr, newUserID)
		}
	}

	// 6. Bump the referrer's joined counter.
	updated, incErr := s.repo.IncrementAccountCounters(ctx, acc.UserID, 1, 0)
	if incErr != nil {
		return s.logHookErr(op, incErr, newUserID)
	}

	// Reward/welcome-coupon rules come from the admin-editable settings.
	eff := s.effective(ctx)

	// 7. Optional welcome coupon for the new user.
	s.grantWelcomeCoupon(ctx, eff, newUserID, phone)

	// 8. Reward the referrer for every full batch of joins.
	s.maybeGrantReward(ctx, eff, updated)

	return nil
}

func (s StorefrontService) grantWelcomeCoupon(ctx context.Context, eff Config, newUserID, phone string) {
	if !eff.WelcomeDiscountEnabled {
		return
	}
	d, err := s.mintCoupon(ctx, newUserID, eff.WelcomeDiscountType, eff.WelcomeDiscountValue, 0, eff.welcomeValidDays())
	if err != nil {
		logger.L().Error("referral: failed to mint welcome coupon", "err", err, "user_id", newUserID)
		return
	}
	if sErr := s.sms.SendPatternSMS(ctx, phone, []string{d.Code}, mellipayamak.ReferralWelcomeKey); sErr != nil {
		logger.L().Error("referral: failed to send welcome coupon sms", "err", sErr, "user_id", newUserID)
	}
}

func (s StorefrontService) maybeGrantReward(ctx context.Context, eff Config, acc entity.ReferralAccount) {
	threshold := eff.rewardThreshold()
	if acc.PendingProgress() < threshold {
		return
	}

	batch, err := s.repo.ListUnrewardedJoined(ctx, acc.UserID, threshold)
	if err != nil {
		logger.L().Error("referral: failed to list unrewarded joins", "err", err, "user_id", acc.UserID)
		return
	}
	if len(batch) < threshold {
		// Counters and rows drifted; don't grant on an incomplete batch.
		logger.L().Error("referral: unrewarded batch smaller than threshold",
			"user_id", acc.UserID, "have", len(batch), "threshold", threshold)
		return
	}

	d, mErr := s.mintCoupon(ctx, acc.UserID, eff.RewardDiscountType, eff.RewardDiscountValue,
		eff.RewardMinimumPurchaseAmount, eff.rewardValidDays())
	if mErr != nil {
		logger.L().Error("referral: failed to mint reward coupon", "err", mErr, "user_id", acc.UserID)
		return
	}

	ids := make([]string, 0, threshold)
	for _, r := range batch {
		ids = append(ids, r.ID)
	}

	reward := entity.ReferralReward{
		ReferrerUserID: acc.UserID,
		DiscountID:     d.ID,
		DiscountCode:   d.Code,
		ReferralIDs:    ids,
	}
	reward.BeforeCreate()
	if cErr := s.repo.CreateReward(ctx, reward); cErr != nil {
		logger.L().Error("referral: failed to persist reward", "err", cErr, "user_id", acc.UserID)
		return
	}
	if rErr := s.repo.MarkReferralsRewarded(ctx, ids, reward.ID); rErr != nil {
		logger.L().Error("referral: failed to stamp rewarded referrals", "err", rErr, "user_id", acc.UserID)
	}
	if _, incErr := s.repo.IncrementAccountCounters(ctx, acc.UserID, 0, threshold); incErr != nil {
		logger.L().Error("referral: failed to bump rewarded_count", "err", incErr, "user_id", acc.UserID)
	}

	referrer, uErr := s.userRepo.GetUserByID(ctx, acc.UserID)
	if uErr != nil {
		logger.L().Error("referral: failed to load referrer for sms", "err", uErr, "user_id", acc.UserID)
		return
	}
	if sErr := s.sms.SendPatternSMS(ctx, referrer.PhoneNumber,
		[]string{strconv.Itoa(threshold), d.Code}, mellipayamak.ReferralRewardKey); sErr != nil {
		logger.L().Error("referral: failed to send reward sms", "err", sErr, "user_id", acc.UserID)
	}
}

func (s StorefrontService) logHookErr(op string, err error, userID string) error {
	logger.L().Error("referral: signup attribution failed", "op", op, "err", err, "user_id", userID)
	return richerror.New(op).WithWarpError(err)
}

func normalizeCode(c string) string {
	return strings.ToUpper(strings.TrimSpace(c))
}

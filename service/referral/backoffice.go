package referralservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// BackofficeService is the admin-panel view of the referral program: a
// leaderboard of referrers, a per-referrer detail, and program-wide stats.
// Read-only.
type BackofficeService struct {
	cfg      Config
	repo     Repository
	userRepo UserProvider
}

func NewBackofficeService(cfg Config, repo Repository, userRepo UserProvider) BackofficeService {
	return BackofficeService{cfg: cfg, repo: repo, userRepo: userRepo}
}

// GetSettings returns the current effective reward/welcome configuration for the
// admin editor. Source is "config" until an admin saves an override.
func (s BackofficeService) GetSettings(ctx context.Context) (referralparam.GetReferralSettingsResponse, error) {
	const op = "referralservice.BackofficeService.GetSettings"

	resp := referralparam.GetReferralSettingsResponse{
		ProgramEnabled: s.cfg.Enabled,
		InviteBaseURL:  s.cfg.InviteBaseURL,
		Source:         "config",
	}

	st, err := s.repo.GetSettings(ctx)
	if err != nil {
		if richerror.GetKind(err) == richerror.KindNotFound {
			resp.Settings = toBody(s.cfg.toSettings())
			return resp, nil
		}
		return referralparam.GetReferralSettingsResponse{}, richerror.New(op).WithWarpError(err)
	}

	resp.Source = "custom"
	resp.Settings = toBody(st)
	updatedAt := st.UpdatedAt
	resp.UpdatedAt = &updatedAt
	resp.UpdatedBy = st.UpdatedBy
	return resp, nil
}

// UpdateSettings replaces the reward/welcome configuration.
func (s BackofficeService) UpdateSettings(ctx context.Context, req referralparam.UpdateReferralSettingsRequest) (referralparam.UpdateReferralSettingsResponse, error) {
	const op = "referralservice.BackofficeService.UpdateSettings"

	st := entity.ReferralSettings{
		RewardThreshold:             req.RewardThreshold,
		RewardDiscountType:          req.RewardDiscountType,
		RewardDiscountValue:         req.RewardDiscountValue,
		RewardValidDays:             req.RewardValidDays,
		RewardMinimumPurchaseAmount: req.RewardMinimumPurchaseAmount,
		WelcomeDiscountEnabled:      req.WelcomeDiscountEnabled,
		WelcomeDiscountType:         req.WelcomeDiscountType,
		WelcomeDiscountValue:        req.WelcomeDiscountValue,
		WelcomeValidDays:            req.WelcomeValidDays,
		UpdatedBy:                   req.AdminUserID,
	}
	if err := s.repo.UpsertSettings(ctx, st); err != nil {
		return referralparam.UpdateReferralSettingsResponse{}, richerror.New(op).WithWarpError(err)
	}
	return referralparam.UpdateReferralSettingsResponse{Settings: toBody(st)}, nil
}

func toBody(st entity.ReferralSettings) referralparam.ReferralSettingsBody {
	return referralparam.ReferralSettingsBody{
		RewardThreshold:             st.RewardThreshold,
		RewardDiscountType:          st.RewardDiscountType,
		RewardDiscountValue:         st.RewardDiscountValue,
		RewardValidDays:             st.RewardValidDays,
		RewardMinimumPurchaseAmount: st.RewardMinimumPurchaseAmount,
		WelcomeDiscountEnabled:      st.WelcomeDiscountEnabled,
		WelcomeDiscountType:         st.WelcomeDiscountType,
		WelcomeDiscountValue:        st.WelcomeDiscountValue,
		WelcomeValidDays:            st.WelcomeValidDays,
	}
}

// List returns the referrer leaderboard (sorted by total joined, desc).
func (s BackofficeService) List(ctx context.Context, req referralparam.AdminListReferralsRequest) (referralparam.AdminListReferralsResponse, error) {
	const op = "referralservice.BackofficeService.List"

	accounts, total, err := s.repo.AllAccounts(ctx, req.MinJoined, req.Pagination)
	if err != nil {
		return referralparam.AdminListReferralsResponse{}, richerror.New(op).WithWarpError(err)
	}

	usersByID := s.resolveUsers(ctx, accountUserIDs(accounts))

	items := make([]referralparam.AdminAccountRow, 0, len(accounts))
	for _, a := range accounts {
		rewardsEarned, cErr := s.repo.CountRewardsByReferrer(ctx, a.UserID)
		if cErr != nil {
			return referralparam.AdminListReferralsResponse{}, richerror.New(op).WithWarpError(cErr)
		}
		items = append(items, s.toRow(a, usersByID[a.UserID], int(rewardsEarned)))
	}

	return referralparam.AdminListReferralsResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(total),
		},
		Items: items,
	}, nil
}

// Get returns one referrer's account, their (paginated) invitees and all their
// earned reward codes.
func (s BackofficeService) Get(ctx context.Context, req referralparam.AdminGetReferralRequest) (referralparam.AdminGetReferralResponse, error) {
	const op = "referralservice.BackofficeService.Get"

	acc, err := s.repo.FindAccountByUserID(ctx, req.UserID)
	if err != nil {
		return referralparam.AdminGetReferralResponse{}, richerror.New(op).WithWarpError(err)
	}

	rows, total, err := s.repo.ListReferralsByReferrer(ctx, req.UserID, req.Pagination)
	if err != nil {
		return referralparam.AdminGetReferralResponse{}, richerror.New(op).WithWarpError(err)
	}

	rewards, err := s.repo.ListRewardsByReferrer(ctx, req.UserID)
	if err != nil {
		return referralparam.AdminGetReferralResponse{}, richerror.New(op).WithWarpError(err)
	}

	// invitee names
	joinedIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.ReferredUserID != "" {
			joinedIDs = append(joinedIDs, r.ReferredUserID)
		}
	}
	inviteeByID := s.resolveUsers(ctx, joinedIDs)
	referrerByID := s.resolveUsers(ctx, []string{req.UserID})

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
		if u, ok := inviteeByID[r.ReferredUserID]; ok {
			item.Name = u.Name
			item.PhoneMasked = maskPhone(u.PhoneNumber)
		} else if r.InvitedPhone != "" {
			item.PhoneMasked = maskPhone(r.InvitedPhone)
		}
		items = append(items, item)
	}

	rewardItems := make([]referralparam.RewardItem, 0, len(rewards))
	for _, rw := range rewards {
		rewardItems = append(rewardItems, referralparam.RewardItem{
			ID:           rw.ID,
			DiscountCode: rw.DiscountCode,
			CreatedAt:    rw.CreatedAt,
		})
	}

	return referralparam.AdminGetReferralResponse{
		Account: s.toRow(acc, referrerByID[req.UserID], len(rewards)),
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(total),
		},
		Referrals: items,
		Rewards:   rewardItems,
	}, nil
}

// Stats returns the program-wide dashboard summary.
func (s BackofficeService) Stats(ctx context.Context) (referralparam.AdminReferralStatsResponse, error) {
	const op = "referralservice.BackofficeService.Stats"

	accounts, invites, joined, rewards, err := s.repo.Stats(ctx)
	if err != nil {
		return referralparam.AdminReferralStatsResponse{}, richerror.New(op).WithWarpError(err)
	}

	var rate float64
	if invites > 0 {
		rate = float64(joined) / float64(invites)
	}

	return referralparam.AdminReferralStatsResponse{
		TotalAccounts:  accounts,
		TotalInvites:   invites,
		TotalJoined:    joined,
		TotalRewards:   rewards,
		ConversionRate: rate,
	}, nil
}

func (s BackofficeService) toRow(a entity.ReferralAccount, u entity.User, rewardsEarned int) referralparam.AdminAccountRow {
	row := referralparam.AdminAccountRow{
		UserID:          a.UserID,
		Code:            a.Code,
		TotalJoined:     a.TotalJoined,
		RewardedCount:   a.RewardedCount,
		PendingProgress: a.PendingProgress(),
		RewardsEarned:   rewardsEarned,
		CreatedAt:       a.CreatedAt,
	}
	if u.ID != "" {
		row.ReferrerName = u.Name
		row.ReferrerPhone = u.PhoneNumber
	}
	return row
}

func (s BackofficeService) resolveUsers(ctx context.Context, ids []string) map[string]entity.User {
	out := make(map[string]entity.User)
	if len(ids) == 0 {
		return out
	}
	users, err := s.userRepo.FindManyByIDs(ctx, ids)
	if err != nil {
		return out
	}
	for _, u := range users {
		out[u.ID] = u
	}
	return out
}

func accountUserIDs(accounts []entity.ReferralAccount) []string {
	ids := make([]string, 0, len(accounts))
	for _, a := range accounts {
		ids = append(ids, a.UserID)
	}
	return ids
}

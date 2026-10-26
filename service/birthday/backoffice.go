package birthdayservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	birthdayparam "github.com/abolfazlnorzad/steve/param/birthday"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// GetSettings returns the effective configuration for the admin editor.
// Source is "config" until an admin saves an override.
func (s Service) GetSettings(ctx context.Context) (birthdayparam.GetSettingsResponse, error) {
	const op = "birthdayservice.GetSettings"

	st, custom, err := s.effective(ctx)
	if err != nil {
		return birthdayparam.GetSettingsResponse{}, richerror.New(op).WithWarpError(err)
	}
	resp := birthdayparam.GetSettingsResponse{Settings: toBody(st), Source: "config"}
	if custom {
		resp.Source = "custom"
		updatedAt := st.UpdatedAt
		resp.UpdatedAt = &updatedAt
		resp.UpdatedBy = st.UpdatedBy
	}
	return resp, nil
}

// UpdateSettings replaces the configuration. It takes effect on the next run.
func (s Service) UpdateSettings(ctx context.Context, req birthdayparam.UpdateSettingsRequest) (birthdayparam.UpdateSettingsResponse, error) {
	const op = "birthdayservice.UpdateSettings"

	st := entity.BirthdaySettings{
		Enabled:               req.Enabled,
		DaysBefore:            req.DaysBefore,
		SendHour:              req.SendHour,
		DiscountType:          req.DiscountType,
		DiscountValue:         req.DiscountValue,
		ValidDays:             req.ValidDays,
		MinimumPurchaseAmount: req.MinimumPurchaseAmount,
		UpdatedBy:             req.AdminUserID,
	}
	if err := s.repo.UpsertSettings(ctx, st); err != nil {
		return birthdayparam.UpdateSettingsResponse{}, richerror.New(op).WithWarpError(err)
	}
	return birthdayparam.UpdateSettingsResponse{Settings: toBody(st)}, nil
}

func (s Service) ListRewards(ctx context.Context, req birthdayparam.ListRewardsRequest) (birthdayparam.ListRewardsResponse, error) {
	const op = "birthdayservice.ListRewards"

	items, total, err := s.repo.ListRewards(ctx, req.Pagination)
	if err != nil {
		return birthdayparam.ListRewardsResponse{}, richerror.New(op).WithWarpError(err)
	}
	if items == nil {
		items = []entity.BirthdayReward{}
	}
	return birthdayparam.ListRewardsResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(total),
		},
		Items: items,
	}, nil
}

func toBody(st entity.BirthdaySettings) birthdayparam.SettingsBody {
	return birthdayparam.SettingsBody{
		Enabled:               st.Enabled,
		DaysBefore:            st.DaysBefore,
		SendHour:              st.SendHour,
		DiscountType:          st.DiscountType,
		DiscountValue:         st.DiscountValue,
		ValidDays:             st.ValidDays,
		MinimumPurchaseAmount: st.MinimumPurchaseAmount,
	}
}

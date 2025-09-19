package settingservice

import (
	"context"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontShippingService struct {
	repo ShippingReader
}

func NewStorefrontShippingService(r ShippingReader) StorefrontShippingService {
	return StorefrontShippingService{
		repo: r,
	}
}

func (s StorefrontShippingService) GetActiveShipping(ctx context.Context, _ settingparam.GetActiveShippingRequest) (settingparam.GetActiveShippingResponse, error) {
	const op = "settingservice.GetActiveShipping"

	t := true
	resp, err := s.repo.GetAll(ctx, settingparam.FilterShipping{IsActive: &t}, commonparam.PaginationRequest{
		PageSize:   100,
		PageNumber: 1,
	})
	if err != nil {
		return settingparam.GetActiveShippingResponse{}, richerror.New(op).WithWarpError(err)
	}

	return settingparam.GetActiveShippingResponse{Shipping: resp.Shipping}, nil
}

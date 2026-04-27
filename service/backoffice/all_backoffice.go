package backofficeservice

import (
	"context"

	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) GetAll(ctx context.Context, req backofficeparam.AllBackofficeRequest) (backofficeparam.AllBackofficeResponse, error) {
	const op = "backofficeservice.GetAll"

	resp, err := s.repo.AllBackoffice(ctx, req.Filter, req.Pagination)
	if err != nil {
		return backofficeparam.AllBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return backofficeparam.AllBackofficeResponse{
		Backoffice: resp.Backoffice,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
	}, nil
}

func (s Service) GetSuperUserBackOffices(ctx context.Context, req backofficeparam.GetSuperUserBackofficeRequest) (backofficeparam.GetSuperUserBackofficeResponse, error) {
	const op = "backofficeservice.GetSuperUserBackOffices"

	resp, err := s.repo.AllSuperUserBackOffices(ctx)
	if err != nil {
		return backofficeparam.GetSuperUserBackofficeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return backofficeparam.GetSuperUserBackofficeResponse{
		Backoffice: resp.Backoffice,
	}, nil
}

package supplierservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	supplierparam "github.com/abolfazlnorzad/steve/param/supplier"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) Create(ctx context.Context, req supplierparam.CreateSupplierRequest) (supplierparam.CreateSupplierResponse, error) {
	const op = "supplierservice.Create"

	sup := entity.Supplier{
		Name:        req.Name,
		ShopName:    req.ShopName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Note:        req.Note,
	}
	sup.BeforeCreate()

	if err := s.repo.Create(ctx, sup); err != nil {
		return supplierparam.CreateSupplierResponse{}, richerror.New(op).WithWarpError(err)
	}
	return supplierparam.CreateSupplierResponse{Supplier: sup}, nil
}

func (s BackofficeService) Get(ctx context.Context, id string) (supplierparam.GetSupplierResponse, error) {
	const op = "supplierservice.Get"

	sup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return supplierparam.GetSupplierResponse{}, richerror.New(op).WithWarpError(err)
	}
	return supplierparam.GetSupplierResponse{Supplier: sup}, nil
}

func (s BackofficeService) List(ctx context.Context, req supplierparam.ListSuppliersRequest) (supplierparam.ListSuppliersResponse, error) {
	const op = "supplierservice.List"

	suppliers, total, err := s.repo.All(ctx, req.Search, req.IsActive, req.Pagination)
	if err != nil {
		return supplierparam.ListSuppliersResponse{}, richerror.New(op).WithWarpError(err)
	}
	return supplierparam.ListSuppliersResponse{
		Suppliers: suppliers,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

func (s BackofficeService) Update(ctx context.Context, req supplierparam.UpdateSupplierRequest) (supplierparam.UpdateSupplierResponse, error) {
	const op = "supplierservice.Update"

	sup, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return supplierparam.UpdateSupplierResponse{}, richerror.New(op).WithWarpError(err)
	}

	if req.Name != nil {
		sup.Name = *req.Name
	}
	if req.ShopName != nil {
		sup.ShopName = *req.ShopName
	}
	if req.PhoneNumber != nil {
		sup.PhoneNumber = *req.PhoneNumber
	}
	if req.Address != nil {
		sup.Address = *req.Address
	}
	if req.Note != nil {
		sup.Note = *req.Note
	}
	if req.IsActive != nil {
		sup.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, sup); err != nil {
		return supplierparam.UpdateSupplierResponse{}, richerror.New(op).WithWarpError(err)
	}
	return supplierparam.UpdateSupplierResponse{Supplier: sup}, nil
}

func (s BackofficeService) Delete(ctx context.Context, id string) error {
	const op = "supplierservice.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

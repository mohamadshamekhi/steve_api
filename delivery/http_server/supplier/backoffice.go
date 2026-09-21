package supplierhandler

import (
	"strconv"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
	supplierparam "github.com/abolfazlnorzad/steve/param/supplier"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// pagination reads page_number/page_size, falling back to the defaults when
// either is absent or unparsable.
func pagination(c echo.Context) commonparam.PaginationRequest {
	pageNumber, err := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if err != nil {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, err := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if err != nil {
		pageSize = commonparam.DefaultPageSize
	}
	return commonparam.PaginationRequest{PageNumber: uint(pageNumber), PageSize: uint(pageSize)}
}

func (h Handler) List(c echo.Context) error {
	req := supplierparam.ListSuppliersRequest{
		Search:     c.QueryParam("search"),
		Pagination: pagination(c),
	}
	if raw := c.QueryParam("is_active"); raw != "" {
		isActive := raw == "true"
		req.IsActive = &isActive
	}

	resp, err := h.backofficeSvc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Get(c echo.Context) error {
	resp, err := h.backofficeSvc.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Create(c echo.Context) error {
	const op = "supplierhandler.Create"

	var req supplierparam.CreateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	if err := h.validator.ValidateCreateRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Update(c echo.Context) error {
	const op = "supplierhandler.Update"

	var req supplierparam.UpdateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.ID = c.Param("id")

	if err := h.validator.ValidateUpdateRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return h.success(c, resp)
}

func (h Handler) Delete(c echo.Context) error {
	if err := h.backofficeSvc.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return h.success(c, nil)
}

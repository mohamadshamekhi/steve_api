package partnerservice

import (
	"context"
	"fmt"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
)

// OrderManager reads and writes partner orders for the fulfilment flow.
// Satisfied by *mongodborder.DB.
type OrderManager interface {
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	SetPartnerHandover(ctx context.Context, orderID, handoverBy string, handoverAt time.Time) error
	AllPartnerOrders(ctx context.Context, status string, p commonparam.PaginationRequest) (orderservice.ListResult, error)
}

// InPersonPaymentRecorder records a counter payment. Satisfied by
// paymentservice.Service.
type InPersonPaymentRecorder interface {
	CreatePaymentForBackofficePanel(ctx context.Context, req paymentparam.CreatePaymentForBackofficePanelRequest) (paymentparam.CreatePaymentForBackofficePanelResponse, error)
}

type BackofficeService struct {
	repo       Repository
	products   ProductProvider
	users      UserProvider
	backoffice BackofficeProvisioner
	orders     OrderManager
	payments   InPersonPaymentRecorder
	sms        SMSSender
	cfg        Config
}

func NewBackofficeService(
	repo Repository,
	products ProductProvider,
	users UserProvider,
	backoffice BackofficeProvisioner,
	orders OrderManager,
	payments InPersonPaymentRecorder,
	sms SMSSender,
	cfg Config,
) BackofficeService {
	return BackofficeService{
		repo:       repo,
		products:   products,
		users:      users,
		backoffice: backoffice,
		orders:     orders,
		payments:   payments,
		sms:        sms,
		cfg:        cfg,
	}
}

// productSoldToPartners is the effective eligibility used in the admin list —
// mirrors ShopService.productEligible.
func (s BackofficeService) productSoldToPartners(p entity.Product) bool {
	if p.PartnerExcluded {
		return false
	}
	return p.PartnerEligible || s.cfg.AllProductsEligible
}

func (s BackofficeService) List(ctx context.Context, req partnerparam.ListRequest) (partnerparam.ListResponse, error) {
	const op = "partnerservice.List"
	partners, total, err := s.repo.All(ctx, req.Status, req.Pagination)
	if err != nil {
		return partnerparam.ListResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.ListResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(total),
		},
		Partners: partners,
	}, nil
}

func (s BackofficeService) Get(ctx context.Context, id string) (partnerparam.GetResponse, error) {
	const op = "partnerservice.Get"
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return partnerparam.GetResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.GetResponse{Partner: p}, nil
}

func (s BackofficeService) Approve(ctx context.Context, id, reviewerID string, req partnerparam.ApproveRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.Approve"
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	if gErr := p.GuardTransition(entity.PartnerStatusApproved); gErr != nil {
		return partnerparam.EmptyResponse{}, gErr
	}

	applicant, aErr := s.resolveApplicant(ctx, p.UserID)
	if aErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(aErr)
	}

	backofficeUserID, pErr := s.provisionBackofficeAccount(ctx, applicant, p.NationalID)
	if pErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(pErr)
	}

	now := time.Now()
	p.Status = entity.PartnerStatusApproved
	p.DiscountPercent = req.DiscountPercent
	p.BackofficeUserID = backofficeUserID
	p.RejectionReason = ""
	p.ReviewedBy = reviewerID
	p.ReviewedAt = &now
	p.UpdatedAt = now
	if uErr := s.repo.Update(ctx, p); uErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	s.notifyPartnerPhone(applicant.PhoneNumber, fmt.Sprintf("درخواست همکاری شما تایید شد. با همین شماره وارد پنل همکاری شوید. درصد تخفیف: %d%%", req.DiscountPercent))
	return partnerparam.EmptyResponse{}, nil
}

// resolveApplicant loads the storefront user who submitted the application.
func (s BackofficeService) resolveApplicant(ctx context.Context, userID string) (entity.User, error) {
	const op = "partnerservice.resolveApplicant"
	if s.users == nil {
		return entity.User{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithMessageKey("global.something_wrong")
	}
	users, err := s.users.FindManyByIDs(ctx, []string{userID})
	if err != nil {
		return entity.User{}, err
	}
	if len(users) == 0 {
		return entity.User{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}
	return users[0], nil
}

// provisionBackofficeAccount find-or-creates the restricted backoffice account
// the partner logs in with, and ensures it has the `partner` role.
func (s BackofficeService) provisionBackofficeAccount(ctx context.Context, applicant entity.User, nationalID string) (string, error) {
	const op = "partnerservice.provisionBackofficeAccount"

	existing, err := s.backoffice.GetBackOfficeUserByPhoneNumber(ctx, applicant.PhoneNumber)
	if err != nil && richerror.GetKind(err) != richerror.KindNotFound {
		return "", richerror.New(op).WithWarpError(err)
	}

	if err == nil {
		if rErr := s.backoffice.AddRoleToBackofficeUser(applicant.PhoneNumber, key.Partner); rErr != nil {
			return "", richerror.New(op).WithWarpError(rErr)
		}
		return existing.ID, nil
	}

	newUser := entity.BackOffice{
		Name:         applicant.Name,
		PhoneNumber:  applicant.PhoneNumber,
		NationalCode: nationalID,
		IsActive:     true,
		Roles:        []string{key.Partner},
	}
	newUser.BeforeCreate()
	if cErr := s.backoffice.AddNewBackofficeUser(ctx, newUser); cErr != nil {
		return "", richerror.New(op).WithWarpError(cErr)
	}
	return newUser.ID, nil
}

func (s BackofficeService) Reject(ctx context.Context, id, reviewerID string, req partnerparam.RejectRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.Reject"
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	if gErr := p.GuardTransition(entity.PartnerStatusRejected); gErr != nil {
		return partnerparam.EmptyResponse{}, gErr
	}
	now := time.Now()
	p.Status = entity.PartnerStatusRejected
	p.RejectionReason = req.Reason
	p.ReviewedBy = reviewerID
	p.ReviewedAt = &now
	p.UpdatedAt = now
	if uErr := s.repo.Update(ctx, p); uErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	if applicant, aErr := s.resolveApplicant(ctx, p.UserID); aErr == nil {
		s.notifyPartnerPhone(applicant.PhoneNumber, "درخواست همکاری شما پذیرفته نشد.")
	}
	return partnerparam.EmptyResponse{}, nil
}

// setStatus is the shared body for suspend / unsuspend.
func (s BackofficeService) setStatus(ctx context.Context, op, id string, to entity.PartnerStatus) (partnerparam.EmptyResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	if gErr := p.GuardTransition(to); gErr != nil {
		return partnerparam.EmptyResponse{}, gErr
	}
	p.Status = to
	p.UpdatedAt = time.Now()
	if uErr := s.repo.Update(ctx, p); uErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	return partnerparam.EmptyResponse{}, nil
}

func (s BackofficeService) Suspend(ctx context.Context, id string) (partnerparam.EmptyResponse, error) {
	return s.setStatus(ctx, "partnerservice.Suspend", id, entity.PartnerStatusSuspended)
}

func (s BackofficeService) Unsuspend(ctx context.Context, id string) (partnerparam.EmptyResponse, error) {
	return s.setStatus(ctx, "partnerservice.Unsuspend", id, entity.PartnerStatusApproved)
}

func (s BackofficeService) UpdatePricing(ctx context.Context, id string, req partnerparam.UpdatePricingRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.UpdatePricing"
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !p.IsApproved() {
		return partnerparam.EmptyResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("partner.must_be_approved_for_pricing")
	}
	p.DiscountPercent = req.DiscountPercent
	p.UpdatedAt = time.Now()
	if uErr := s.repo.Update(ctx, p); uErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	return partnerparam.EmptyResponse{}, nil
}

func (s BackofficeService) SetProductConfig(ctx context.Context, productID string, req partnerparam.SetProductConfigRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.SetProductConfig"
	if _, err := s.products.FindByID(ctx, productID); err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	percent := req.DiscountPercent
	if percent != nil && *percent > 100 {
		return partnerparam.EmptyResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).WithMessageKey("partner.invalid_discount_percent")
	}
	if err := s.products.SetPartnerConfig(ctx, productID, req.Eligible, percent); err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.EmptyResponse{}, nil
}

func (s BackofficeService) BulkSetEligible(ctx context.Context, req partnerparam.BulkSetEligibleRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.BulkSetEligible"
	if len(req.ProductIDs) == 0 {
		return partnerparam.EmptyResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).WithMessageKey("global.validation_failed")
	}
	if err := s.products.BulkSetPartnerEligible(ctx, req.ProductIDs, req.Eligible); err != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(err)
	}
	return partnerparam.EmptyResponse{}, nil
}

// RecordInPersonPayment records a counter payment for a partner order and moves
// it into processing (ready for handover).
func (s BackofficeService) RecordInPersonPayment(ctx context.Context, orderID string, req partnerparam.RecordInPersonPaymentRequest) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.RecordInPersonPayment"
	order, err := s.partnerOrder(ctx, op, orderID)
	if err != nil {
		return partnerparam.EmptyResponse{}, err
	}
	if order.Status != entity.OrderStatusPendingPayment {
		return partnerparam.EmptyResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("partner.order_not_pending_payment")
	}
	if _, pErr := s.payments.CreatePaymentForBackofficePanel(ctx, paymentparam.CreatePaymentForBackofficePanelRequest{
		OrderID:     orderID,
		Description: req.Description,
	}); pErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(pErr)
	}
	if uErr := s.orders.UpdateOrderStatus(ctx, orderID, entity.OrderStatusProcessing); uErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(uErr)
	}
	return partnerparam.EmptyResponse{}, nil
}

// RecordHandover marks a paid partner order as handed over in person and
// delivered.
func (s BackofficeService) RecordHandover(ctx context.Context, orderID, handlerID string) (partnerparam.EmptyResponse, error) {
	const op = "partnerservice.RecordHandover"
	order, err := s.partnerOrder(ctx, op, orderID)
	if err != nil {
		return partnerparam.EmptyResponse{}, err
	}
	if order.Status != entity.OrderStatusProcessing {
		return partnerparam.EmptyResponse{}, richerror.New(op).
			WithKind(richerror.KindBadRequest).WithMessageKey("partner.order_not_ready_for_handover")
	}
	if hErr := s.orders.SetPartnerHandover(ctx, orderID, handlerID, time.Now()); hErr != nil {
		return partnerparam.EmptyResponse{}, richerror.New(op).WithWarpError(hErr)
	}
	return partnerparam.EmptyResponse{}, nil
}

// ListProducts powers the admin's "partner product config" screen: every
// product with its current eligibility + override percent.
func (s BackofficeService) ListProducts(ctx context.Context, req partnerparam.ProductListRequest) (partnerparam.ProductListResponse, error) {
	const op = "partnerservice.ListProducts"

	filter := commonparam.FilterRequest{}
	if req.Name != "" {
		filter.Filter = append(filter.Filter, map[string]any{"name": req.Name})
	}
	if req.BrandID != "" {
		filter.Filter = append(filter.Filter, map[string]any{"brand_id": req.BrandID})
	}
	if req.CategoryID != "" {
		filter.Filter = append(filter.Filter, map[string]any{"category_id": req.CategoryID})
	}
	// eligible=false → only products the admin explicitly excluded.
	if req.Eligible != nil && !*req.Eligible {
		filter.Filter = append(filter.Filter, map[string]any{"partner_excluded": true})
	}

	list, err := s.products.All(ctx, filter, req.Pagination)
	if err != nil {
		return partnerparam.ProductListResponse{}, richerror.New(op).WithWarpError(err)
	}

	items := make([]partnerparam.ProductListItem, 0, len(list.Products))
	for _, p := range list.Products {
		it := partnerparam.ProductListItem{
			ID:                     p.ID,
			Name:                   p.Name,
			Slug:                   p.Slug,
			IsPublished:            p.IsPublished,
			SoldToPartners:         s.productSoldToPartners(p),
			PartnerExcluded:        p.PartnerExcluded,
			PartnerDiscountPercent: p.PartnerDiscountPercent,
			Variants:               make([]partnerparam.ProductListVariant, 0, len(p.Variants)),
		}
		for _, v := range p.Variants {
			it.Variants = append(it.Variants, partnerparam.ProductListVariant{ID: v.ID, Price: v.Price, Quantity: v.Quantity})
		}
		items = append(items, it)
	}

	return partnerparam.ProductListResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(list.TotalCount),
		},
		Products: items,
	}, nil
}

// ListOrders powers the admin's "partner orders" screen.
func (s BackofficeService) ListOrders(ctx context.Context, req partnerparam.OrderListRequest) (partnerparam.OrderListResponse, error) {
	const op = "partnerservice.ListOrders"
	list, err := s.orders.AllPartnerOrders(ctx, req.Status, req.Pagination)
	if err != nil {
		return partnerparam.OrderListResponse{}, richerror.New(op).WithWarpError(err)
	}
	orders := make([]partnerparam.OrderListItem, 0, len(list.Orders))
	for _, o := range list.Orders {
		orders = append(orders, partnerparam.OrderListItem{
			ID:               o.ID,
			PartnerID:        o.PartnerID,
			UserID:           o.UserID,
			Status:           string(o.Status),
			PurchaseType:     string(o.PurchaseType),
			SubtotalAmount:   o.SubtotalAmount,
			FinalAmount:      o.FinalAmount,
			AppliedDiscounts: o.AppliedDiscounts,
			Items:            o.Items,
			HandoverBy:       o.HandoverBy,
			HandoverAt:       o.HandoverAt,
			CreatedAt:        o.CreatedAt,
		})
	}
	return partnerparam.OrderListResponse{
		Pagination: commonparam.PaginationResponse{
			PageSize:   req.Pagination.GetPageSize(),
			PageNumber: req.Pagination.GetPageNumber(),
			Total:      uint(list.TotalCount),
		},
		Orders: orders,
	}, nil
}

func (s BackofficeService) partnerOrder(ctx context.Context, op, orderID string) (entity.Order, error) {
	order, err := s.orders.GetOrderByID(ctx, orderID)
	if err != nil {
		return entity.Order{}, richerror.New(op).WithWarpError(err)
	}
	if !order.IsPartnerOrder {
		return entity.Order{}, richerror.New(op).
			WithKind(richerror.KindForbidden).WithMessageKey("partner.not_a_partner_order")
	}
	return order, nil
}

func (s BackofficeService) notifyPartnerPhone(phone, message string) {
	if s.sms == nil || phone == "" {
		return
	}
	go func() {
		sCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.sms.SendSimple(sCtx, phone, message); err != nil {
			logger.L().Error("partner: failed to send sms", "err", err)
		}
	}()
}

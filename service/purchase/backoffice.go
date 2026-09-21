package purchaseservice

import (
	"context"
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	purchaseparam "github.com/abolfazlnorzad/steve/param/purchase"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mongodbpurchase "github.com/abolfazlnorzad/steve/repository/mongo/purchase"
)

// buildItems resolves every requested line against the catalog, so the stored
// invoice carries readable names and can never reference a missing variant.
func (s BackofficeService) buildItems(ctx context.Context, reqItems []purchaseparam.CreatePurchaseItemRequest) ([]entity.PurchaseItem, error) {
	const op = "purchaseservice.buildItems"

	items := make([]entity.PurchaseItem, 0, len(reqItems))
	for _, ri := range reqItems {
		product, err := s.productProvider.FindByID(ctx, ri.ProductID)
		if err != nil {
			return nil, richerror.New(op).WithWarpError(err)
		}
		if product.IsService() {
			return nil, richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("purchase.service_cannot_be_purchased").
				WithMeta(map[string]any{"product_id": ri.ProductID})
		}

		var found bool
		item := entity.PurchaseItem{
			ProductID:   ri.ProductID,
			VariantID:   ri.VariantID,
			ProductName: product.Name,
			Quantity:    ri.Quantity,
			UnitCost:    ri.UnitCost,
			AddToStock:  ri.AddToStock == nil || *ri.AddToStock,
		}
		for _, v := range product.Variants {
			if v.ID == ri.VariantID {
				item.VariantAttributes = v.VariantAttribute
				found = true
			}
		}
		if !found {
			return nil, richerror.New(op).
				WithKind(richerror.KindNotFound).
				WithMessageKey("product.variant_not_found").
				WithMeta(map[string]any{"product_id": ri.ProductID, "variant_id": ri.VariantID})
		}
		items = append(items, item)
	}
	return items, nil
}

func (s BackofficeService) Create(ctx context.Context, req purchaseparam.CreatePurchaseRequest, createdBy string) (purchaseparam.CreatePurchaseResponse, error) {
	const op = "purchaseservice.Create"

	supplier, err := s.supplierProvider.FindByID(ctx, req.SupplierID)
	if err != nil {
		return purchaseparam.CreatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}

	items, err := s.buildItems(ctx, req.Items)
	if err != nil {
		return purchaseparam.CreatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}

	p := entity.Purchase{
		SupplierID:    supplier.ID,
		SupplierName:  supplier.Name,
		InvoiceNumber: req.InvoiceNumber,
		Items:         items,
		PaidAmount:    req.PaidAmount,
		Note:          req.Note,
		OrderID:       req.OrderID,
		CreatedBy:     createdBy,
	}
	if req.PurchasedAt != nil {
		p.PurchasedAt = *req.PurchasedAt
	}
	p.BeforeCreate()

	if req.Confirm {
		// Apply to stock first: if the ledger write fails, nothing is stored
		// and the admin can retry without a phantom confirmed invoice.
		if err := s.inventory.ApplyPurchase(ctx, p); err != nil {
			return purchaseparam.CreatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
		}
		now := time.Now()
		p.Status = entity.PurchaseStatusConfirmed
		p.ConfirmedAt = &now
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return purchaseparam.CreatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.CreatePurchaseResponse{Purchase: p}, nil
}

func (s BackofficeService) Get(ctx context.Context, id string) (purchaseparam.GetPurchaseResponse, error) {
	const op = "purchaseservice.Get"

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return purchaseparam.GetPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.GetPurchaseResponse{Purchase: p}, nil
}

func (s BackofficeService) List(ctx context.Context, req purchaseparam.ListPurchasesRequest) (purchaseparam.ListPurchasesResponse, error) {
	const op = "purchaseservice.List"

	filter := mongodbpurchase.PurchaseFilter{
		SupplierID:   req.SupplierID,
		Status:       req.Status,
		PaymentState: req.PaymentState,
		OrderID:      req.OrderID,
		From:         req.From,
		To:           req.To,
	}

	purchases, total, err := s.repo.All(ctx, filter, req.Pagination)
	if err != nil {
		return purchaseparam.ListPurchasesResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.ListPurchasesResponse{
		Purchases: purchases,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

// Update edits a draft invoice. A confirmed invoice is immutable: its cost
// has already been folded into the weighted average, so changing it after the
// fact would silently rewrite history. Cancel and re-enter instead.
func (s BackofficeService) Update(ctx context.Context, req purchaseparam.UpdatePurchaseRequest) (purchaseparam.UpdatePurchaseResponse, error) {
	const op = "purchaseservice.Update"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return purchaseparam.UpdatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	if p.Status != entity.PurchaseStatusDraft {
		return purchaseparam.UpdatePurchaseResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("purchase.only_draft_editable")
	}

	if req.SupplierID != nil {
		supplier, sErr := s.supplierProvider.FindByID(ctx, *req.SupplierID)
		if sErr != nil {
			return purchaseparam.UpdatePurchaseResponse{}, richerror.New(op).WithWarpError(sErr)
		}
		p.SupplierID = supplier.ID
		p.SupplierName = supplier.Name
	}
	if req.InvoiceNumber != nil {
		p.InvoiceNumber = *req.InvoiceNumber
	}
	if req.Note != nil {
		p.Note = *req.Note
	}
	if req.PaidAmount != nil {
		p.PaidAmount = *req.PaidAmount
	}
	if req.PurchasedAt != nil {
		p.PurchasedAt = *req.PurchasedAt
	}
	if req.Items != nil {
		items, iErr := s.buildItems(ctx, req.Items)
		if iErr != nil {
			return purchaseparam.UpdatePurchaseResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		p.Items = items
	}
	p.Recalculate()

	if err := s.repo.Update(ctx, p); err != nil {
		return purchaseparam.UpdatePurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.UpdatePurchaseResponse{Purchase: p}, nil
}

// Confirm applies a draft invoice to the warehouse and the ledger.
func (s BackofficeService) Confirm(ctx context.Context, id string) (purchaseparam.ConfirmPurchaseResponse, error) {
	const op = "purchaseservice.Confirm"

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return purchaseparam.ConfirmPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	if p.Status != entity.PurchaseStatusDraft {
		return purchaseparam.ConfirmPurchaseResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("purchase.only_draft_confirmable")
	}
	if len(p.Items) == 0 {
		return purchaseparam.ConfirmPurchaseResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("purchase.items_required")
	}

	if err := s.inventory.ApplyPurchase(ctx, p); err != nil {
		return purchaseparam.ConfirmPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}

	now := time.Now()
	p.Status = entity.PurchaseStatusConfirmed
	p.ConfirmedAt = &now
	if err := s.repo.Update(ctx, p); err != nil {
		return purchaseparam.ConfirmPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.ConfirmPurchaseResponse{Purchase: p}, nil
}

// Cancel voids an invoice. A confirmed one also has its stock effect undone.
func (s BackofficeService) Cancel(ctx context.Context, id string) (purchaseparam.CancelPurchaseResponse, error) {
	const op = "purchaseservice.Cancel"

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return purchaseparam.CancelPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	if p.Status == entity.PurchaseStatusCanceled {
		return purchaseparam.CancelPurchaseResponse{}, richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("purchase.already_canceled")
	}

	if p.Status == entity.PurchaseStatusConfirmed {
		if err := s.inventory.ReversePurchase(ctx, p); err != nil {
			return purchaseparam.CancelPurchaseResponse{}, richerror.New(op).WithWarpError(err)
		}
	}

	p.Status = entity.PurchaseStatusCanceled
	if err := s.repo.Update(ctx, p); err != nil {
		return purchaseparam.CancelPurchaseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return purchaseparam.CancelPurchaseResponse{Purchase: p}, nil
}

// SetPaidAmount records a payment against the invoice (settling a supplier's
// account), independently of the stock lifecycle.
func (s BackofficeService) SetPaidAmount(ctx context.Context, req purchaseparam.SetPaidAmountRequest) (purchaseparam.SetPaidAmountResponse, error) {
	const op = "purchaseservice.SetPaidAmount"

	p, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return purchaseparam.SetPaidAmountResponse{}, richerror.New(op).WithWarpError(err)
	}

	previousPaid := p.PaidAmount
	p.PaidAmount = req.PaidAmount
	p.Recalculate()

	if err := s.repo.Update(ctx, p); err != nil {
		return purchaseparam.SetPaidAmountResponse{}, richerror.New(op).WithWarpError(err)
	}

	// Only the newly paid part reaches the cashbox: PaidAmount is a running
	// total, so booking it whole on every call would count earlier
	// instalments again. A correction downwards is left to a manual entry.
	if p.PaidAmount > previousPaid {
		s.recordPayment(ctx, p, p.PaidAmount-previousPaid, req.PaymentMethod)
	}
	return purchaseparam.SetPaidAmountResponse{Purchase: p}, nil
}

// recordPayment books a supplier settlement in the cashbox. A failure here
// must not undo the payment record on the invoice.
func (s BackofficeService) recordPayment(ctx context.Context, p entity.Purchase, amount uint64, method entity.PaymentMethod) {
	if s.cash == nil || amount == 0 {
		return
	}

	t := entity.CashTransaction{
		Direction:   entity.CashOut,
		Method:      method,
		Amount:      amount,
		Description: "پرداخت به تأمین‌کننده بابت فاکتور خرید",
		RefType:     entity.CashRefPurchase,
		RefID:       p.ID,
		SupplierID:  p.SupplierID,
	}
	if err := s.cash.Record(ctx, t); err != nil {
		logger.L().Error("failed to record supplier payment in cashbox", "purchase_id", p.ID, "err", err)
	}
}

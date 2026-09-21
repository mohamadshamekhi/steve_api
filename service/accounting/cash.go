package accountingservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// CreateCashTransaction records a manual till movement (برداشت شخصی، واریز
// اولیه، دریافت نسیه‌ی قبلی، ...).
func (s BackofficeService) CreateCashTransaction(ctx context.Context, req accountingparam.CreateCashTransactionRequest, createdBy string) (accountingparam.CreateCashTransactionResponse, error) {
	const op = "accountingservice.CreateCashTransaction"

	t := entity.CashTransaction{
		Direction:   req.Direction,
		Method:      req.Method,
		Amount:      req.Amount,
		Description: req.Description,
		SupplierID:  req.SupplierID,
		RefType:     req.RefType,
		RefID:       req.RefID,
		CreatedBy:   createdBy,
	}
	if t.RefType == "" {
		t.RefType = entity.CashRefManual
	}
	if req.OccurredAt != nil {
		t.OccurredAt = *req.OccurredAt
	}
	t.BeforeCreate()

	if err := s.cashRepo.Create(ctx, t); err != nil {
		return accountingparam.CreateCashTransactionResponse{}, richerror.New(op).WithWarpError(err)
	}
	return accountingparam.CreateCashTransactionResponse{Transaction: t}, nil
}

// Record is the port other domains use to push an automatic cash line
// (an order paid at the counter, a supplier invoice settled, a refund).
func (s BackofficeService) Record(ctx context.Context, t entity.CashTransaction) error {
	const op = "accountingservice.Record"

	if t.Amount == 0 {
		return nil
	}
	t.IsAuto = true
	t.BeforeCreate()

	if err := s.cashRepo.Create(ctx, t); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

func (s BackofficeService) ListCashTransactions(ctx context.Context, req accountingparam.ListCashTransactionsRequest) (accountingparam.ListCashTransactionsResponse, error) {
	const op = "accountingservice.ListCashTransactions"

	transactions, total, err := s.cashRepo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return accountingparam.ListCashTransactionsResponse{}, richerror.New(op).WithWarpError(err)
	}

	summary, err := s.cashRepo.Summary(ctx, req.Filter)
	if err != nil {
		return accountingparam.ListCashTransactionsResponse{}, richerror.New(op).WithWarpError(err)
	}

	return accountingparam.ListCashTransactionsResponse{
		Transactions: transactions,
		Summary:      summary,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

// CashReport is the صندوق report: the window's totals, split by payment
// method and by day.
func (s BackofficeService) CashReport(ctx context.Context, req accountingparam.CashReportRequest) (accountingparam.CashReportResponse, error) {
	const op = "accountingservice.CashReport"

	summary, err := s.cashRepo.Summary(ctx, req.Filter)
	if err != nil {
		return accountingparam.CashReportResponse{}, richerror.New(op).WithWarpError(err)
	}
	byMethod, err := s.cashRepo.SummaryByMethod(ctx, req.Filter)
	if err != nil {
		return accountingparam.CashReportResponse{}, richerror.New(op).WithWarpError(err)
	}
	byDay, err := s.cashRepo.SummaryByDay(ctx, req.Filter)
	if err != nil {
		return accountingparam.CashReportResponse{}, richerror.New(op).WithWarpError(err)
	}

	return accountingparam.CashReportResponse{
		Summary:  summary,
		ByMethod: byMethod,
		ByDay:    byDay,
	}, nil
}

// DeleteCashTransaction removes a manual line. Auto-recorded lines are owned
// by the document that produced them (order, purchase, expense, return), so
// deleting them here would leave that document claiming a payment the cashbox
// no longer shows.
func (s BackofficeService) DeleteCashTransaction(ctx context.Context, id string) error {
	const op = "accountingservice.DeleteCashTransaction"

	isAuto, err := s.cashRepo.FindByID(ctx, id)
	if err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	if isAuto {
		return richerror.New(op).
			WithKind(richerror.KindForbidden).
			WithMessageKey("cash.auto_transaction_not_deletable")
	}
	if err := s.cashRepo.Delete(ctx, id); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

// ExportCashTransactions returns the filtered rows for a spreadsheet download.
func (s BackofficeService) ExportCashTransactions(ctx context.Context, f accountingparam.CashFilter) ([]entity.CashTransaction, error) {
	const op = "accountingservice.ExportCashTransactions"

	transactions, err := s.cashRepo.AllForExport(ctx, f)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	return transactions, nil
}

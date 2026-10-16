package accountingservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

// CreateExpense books an operating cost. When it is marked paid, the cashbox
// gets a matching payout line so «دخل و خرج» stays complete without the admin
// having to type the same number twice.
func (s BackofficeService) CreateExpense(ctx context.Context, req accountingparam.CreateExpenseRequest, createdBy string) (accountingparam.CreateExpenseResponse, error) {
	const op = "accountingservice.CreateExpense"

	e := entity.Expense{
		Title:         req.Title,
		Category:      req.Category,
		Amount:        req.Amount,
		SupplierID:    req.SupplierID,
		PaymentMethod: req.PaymentMethod,
		IsPaid:        req.IsPaid == nil || *req.IsPaid,
		Note:          req.Note,
		Attachment:    req.Attachment,
		CreatedBy:     createdBy,
	}
	if req.SpentAt != nil {
		e.SpentAt = *req.SpentAt
	}
	e.BeforeCreate()

	if err := s.expenseRepo.Create(ctx, e); err != nil {
		return accountingparam.CreateExpenseResponse{}, richerror.New(op).WithWarpError(err)
	}

	if e.IsPaid {
		s.recordExpensePayment(ctx, e, createdBy)
	}
	return accountingparam.CreateExpenseResponse{Expense: e}, nil
}

// recordExpensePayment writes the cashbox side of a paid expense. A failure
// here must not undo the expense itself — the cost is a fact, the cash line is
// bookkeeping that can be added by hand.
func (s BackofficeService) recordExpensePayment(ctx context.Context, e entity.Expense, createdBy string) {
	t := entity.CashTransaction{
		Direction:   entity.CashOut,
		Method:      e.PaymentMethod,
		Amount:      e.Amount,
		Description: e.Title,
		RefType:     entity.CashRefExpense,
		RefID:       e.ID,
		SupplierID:  e.SupplierID,
		IsAuto:      true,
		OccurredAt:  e.SpentAt,
		CreatedBy:   createdBy,
	}
	t.BeforeCreate()

	if err := s.cashRepo.Create(ctx, t); err != nil {
		logger.L().Error("failed to record cash transaction for expense", "expense_id", e.ID, "err", err)
	}
}

func (s BackofficeService) GetExpense(ctx context.Context, id string) (accountingparam.GetExpenseResponse, error) {
	const op = "accountingservice.GetExpense"

	e, err := s.expenseRepo.FindByID(ctx, id)
	if err != nil {
		return accountingparam.GetExpenseResponse{}, richerror.New(op).WithWarpError(err)
	}
	return accountingparam.GetExpenseResponse{Expense: e}, nil
}

func (s BackofficeService) ListExpenses(ctx context.Context, req accountingparam.ListExpensesRequest) (accountingparam.ListExpensesResponse, error) {
	const op = "accountingservice.ListExpenses"

	expenses, total, err := s.expenseRepo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return accountingparam.ListExpensesResponse{}, richerror.New(op).WithWarpError(err)
	}

	// The header shows the total of the whole filtered set, not of the page:
	// «این ماه چقدر خرج کردم» must not change when you turn the page.
	totalAmount, err := s.expenseRepo.Total(ctx, req.Filter)
	if err != nil {
		return accountingparam.ListExpensesResponse{}, richerror.New(op).WithWarpError(err)
	}

	return accountingparam.ListExpensesResponse{
		Expenses:    expenses,
		TotalAmount: totalAmount,
		Pagination: commonparam.PaginationResponse{
			PageNumber: req.Pagination.GetPageNumber(),
			PageSize:   req.Pagination.GetPageSize(),
			Total:      uint(total),
		},
	}, nil
}

// UpdateExpense edits a recorded cost. Changing the amount of an already-paid
// expense does not rewrite its cashbox line: the till movement happened as it
// happened, so the correction is left to a manual cash entry.
func (s BackofficeService) UpdateExpense(ctx context.Context, req accountingparam.UpdateExpenseRequest, createdBy string) (accountingparam.UpdateExpenseResponse, error) {
	const op = "accountingservice.UpdateExpense"

	e, err := s.expenseRepo.FindByID(ctx, req.ID)
	if err != nil {
		return accountingparam.UpdateExpenseResponse{}, richerror.New(op).WithWarpError(err)
	}
	wasPaid := e.IsPaid

	if req.Title != nil {
		e.Title = *req.Title
	}
	if req.Category != nil {
		e.Category = *req.Category
	}
	if req.Amount != nil {
		e.Amount = *req.Amount
	}
	if req.SupplierID != nil {
		e.SupplierID = *req.SupplierID
	}
	if req.PaymentMethod != nil {
		e.PaymentMethod = *req.PaymentMethod
	}
	if req.IsPaid != nil {
		e.IsPaid = *req.IsPaid
	}
	if req.SpentAt != nil {
		e.SpentAt = *req.SpentAt
	}
	if req.Note != nil {
		e.Note = *req.Note
	}
	if req.Attachment != nil {
		e.Attachment = *req.Attachment
	}

	if err := s.expenseRepo.Update(ctx, e); err != nil {
		return accountingparam.UpdateExpenseResponse{}, richerror.New(op).WithWarpError(err)
	}

	// Only the unpaid → paid transition creates a cash line; anything else
	// would double-count the payout.
	if !wasPaid && e.IsPaid {
		s.recordExpensePayment(ctx, e, createdBy)
	}
	return accountingparam.UpdateExpenseResponse{Expense: e}, nil
}

func (s BackofficeService) DeleteExpense(ctx context.Context, id string) error {
	const op = "accountingservice.DeleteExpense"

	if err := s.expenseRepo.Delete(ctx, id); err != nil {
		return richerror.New(op).WithWarpError(err)
	}
	return nil
}

// ExportExpenses returns the filtered rows for a spreadsheet download.
func (s BackofficeService) ExportExpenses(ctx context.Context, f accountingparam.ExpenseFilter) ([]entity.Expense, error) {
	const op = "accountingservice.ExportExpenses"

	expenses, err := s.expenseRepo.AllForExport(ctx, f)
	if err != nil {
		return nil, richerror.New(op).WithWarpError(err)
	}
	return expenses, nil
}

package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// ExpenseCategory groups operating costs (هزینه‌های عملیاتی) so the profit and
// loss report can show where the money went, not just how much.
type ExpenseCategory string

const (
	ExpenseCategoryRent      ExpenseCategory = "rent"      // اجاره
	ExpenseCategorySalary    ExpenseCategory = "salary"    // حقوق و دستمزد
	ExpenseCategoryUtilities ExpenseCategory = "utilities" // قبوض (برق/آب/گاز/اینترنت)
	ExpenseCategoryShipping  ExpenseCategory = "shipping"  // حمل و نقل
	ExpenseCategoryMarketing ExpenseCategory = "marketing" // تبلیغات
	ExpenseCategoryTax       ExpenseCategory = "tax"       // مالیات و عوارض
	ExpenseCategoryFee       ExpenseCategory = "fee"       // کارمزد درگاه/بانک
	ExpenseCategorySupplies  ExpenseCategory = "supplies"  // ملزومات مصرفی
	ExpenseCategoryRepair    ExpenseCategory = "repair"    // تعمیر و نگهداری
	ExpenseCategoryOther     ExpenseCategory = "other"     // سایر
)

// IsValid reports whether c is one of the known categories.
func (c ExpenseCategory) IsValid() bool {
	switch c {
	case ExpenseCategoryRent, ExpenseCategorySalary, ExpenseCategoryUtilities,
		ExpenseCategoryShipping, ExpenseCategoryMarketing, ExpenseCategoryTax,
		ExpenseCategoryFee, ExpenseCategorySupplies, ExpenseCategoryRepair,
		ExpenseCategoryOther:
		return true
	}
	return false
}

// Expense is one operating cost paid by the shop. Cost of goods is *not* an
// expense here — it is already carried by the stock ledger as COGS, and
// double-counting it would understate profit.
type Expense struct {
	ID       string          `bson:"_id" json:"id"`
	Title    string          `bson:"title" json:"title"`
	Category ExpenseCategory `bson:"category" json:"category"`
	Amount   uint64          `bson:"amount" json:"amount"`
	// SupplierID is set when the expense is owed to a known supplier
	// (a courier, a landlord kept in the suppliers list, ...).
	SupplierID string `bson:"supplier_id,omitempty" json:"supplier_id,omitempty"`
	// PaymentMethod is empty for an unpaid (accrued) expense.
	PaymentMethod PaymentMethod `bson:"payment_method,omitempty" json:"payment_method,omitempty"`
	IsPaid        bool          `bson:"is_paid" json:"is_paid"`
	// SpentAt is the date the expense belongs to for reporting, which may
	// differ from when it was typed in.
	SpentAt    time.Time `bson:"spent_at" json:"spent_at"`
	Note       string    `bson:"note,omitempty" json:"note,omitempty"`
	Attachment string    `bson:"attachment,omitempty" json:"attachment,omitempty"`
	CreatedBy  string    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (e *Expense) BeforeCreate() {
	e.ID = ulid.Make().String()
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	if e.SpentAt.IsZero() {
		e.SpentAt = now
	}
	if e.Category == "" {
		e.Category = ExpenseCategoryOther
	}
}

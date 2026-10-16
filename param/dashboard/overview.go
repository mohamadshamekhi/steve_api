package dashboardparam

import (
	"time"

	"github.com/abolfazlnorzad/steve/entity"
	accountingparam "github.com/abolfazlnorzad/steve/param/accounting"
)

// OverviewRequest is the reporting window, already resolved to Tehran-local
// calendar days by the handler. FromRaw/ToRaw are the YYYY-MM-DD strings
// (empty = default window).
type OverviewRequest struct {
	FromRaw string
	ToRaw   string
	// BackofficeUserID decides which sections are filled.
	BackofficeUserID string
}

type Period struct {
	From     string `json:"from"`
	To       string `json:"to"`
	PrevFrom string `json:"prev_from"`
	PrevTo   string `json:"prev_to"`
}

type Sales struct {
	NetRevenue         uint64  `json:"net_revenue"`
	NetRevenuePrev     uint64  `json:"net_revenue_prev"`
	NetProfit          int64   `json:"net_profit"`
	NetProfitPrev      int64   `json:"net_profit_prev"`
	GrossProfit        int64   `json:"gross_profit"`
	Revenue            uint64  `json:"revenue"`
	Returns            uint64  `json:"returns"`
	NetCOGS            uint64  `json:"net_cogs"`
	Expenses           uint64  `json:"expenses"`
	GrossMarginPercent float64 `json:"gross_margin_percent"`
	NetMarginPercent   float64 `json:"net_margin_percent"`
	UnitsSold          uint64  `json:"units_sold"`
	OrderCount         int64   `json:"order_count"`
	OrderCountPrev     int64   `json:"order_count_prev"`
	AvgOrderValue      uint64  `json:"avg_order_value"`
}

type TrendRow struct {
	Day         string `json:"day"`
	Revenue     uint64 `json:"revenue"`
	GrossProfit int64  `json:"gross_profit"`
	OrderCount  int64  `json:"order_count"`
}

type OrderStatusCount struct {
	Status entity.OrderStatus `bson:"_id" json:"status"`
	Count  int64              `bson:"count" json:"count"`
}

type TopProduct struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	UnitsSold   int64  `json:"units_sold"`
	Revenue     uint64 `json:"revenue"`
	GrossProfit int64  `json:"gross_profit"`
}

type CashDay struct {
	Day      string `json:"day"`
	TotalIn  uint64 `json:"total_in"`
	TotalOut uint64 `json:"total_out"`
}

type CashMethod struct {
	Method  entity.PaymentMethod `json:"method"`
	TotalIn uint64               `json:"total_in"`
}

type Cash struct {
	TotalIn  uint64       `json:"total_in"`
	TotalOut uint64       `json:"total_out"`
	Balance  int64        `json:"balance"`
	ByDay    []CashDay    `json:"by_day"`
	ByMethod []CashMethod `json:"by_method"`
}

type LowStockItem struct {
	ProductID   string                    `json:"product_id"`
	VariantID   string                    `json:"variant_id"`
	ProductName string                    `json:"product_name"`
	Attributes  []entity.VariantAttribute `json:"attributes"`
	Quantity    uint64                    `json:"quantity"`
}

type Inventory struct {
	TotalStockValue   uint64         `json:"total_stock_value"`
	TotalUnits        uint64         `json:"total_units"`
	LowStockCount     uint           `json:"low_stock_count"`
	OutOfStockCount   uint           `json:"out_of_stock_count"`
	LowStockThreshold uint64         `json:"low_stock_threshold"`
	LowStock          []LowStockItem `json:"low_stock"`
}

type SupplierDebtor struct {
	SupplierID   string `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
	Balance      int64  `json:"balance"`
}

type Suppliers struct {
	TotalDebt  int64            `json:"total_debt"`
	TopDebtors []SupplierDebtor `json:"top_debtors"`
}

type RecentOrder struct {
	ID              string             `json:"id"`
	UserName        string             `json:"user_name"`
	UserPhoneNumber string             `json:"user_phone_number"`
	FinalAmount     uint               `json:"final_amount"`
	Status          entity.OrderStatus `json:"status"`
	CreatedAt       time.Time          `json:"created_at"`
}

type Customers struct {
	NewCount     int64 `json:"new_count"`
	NewCountPrev int64 `json:"new_count_prev"`
}

type Repairs struct {
	OpenCount int64 `json:"open_count"`
}

// OverviewResponse is the backoffice home page. A section the caller has no
// permission for is null.
type OverviewResponse struct {
	Period           Period                               `json:"period"`
	Sales            *Sales                               `json:"sales"`
	Trend            []TrendRow                           `json:"trend"`
	OrdersByStatus   []OrderStatusCount                   `json:"orders_by_status"`
	TopProducts      []TopProduct                         `json:"top_products"`
	ExpenseBreakdown []accountingparam.ExpenseCategoryRow `json:"expense_breakdown"`
	Cash             *Cash                                `json:"cash"`
	Inventory        *Inventory                           `json:"inventory"`
	Suppliers        *Suppliers                           `json:"suppliers"`
	RecentOrders     []RecentOrder                        `json:"recent_orders"`
	Customers        *Customers                           `json:"customers"`
	Repairs          *Repairs                             `json:"repairs"`
}

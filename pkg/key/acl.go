package key

const (
	SuperAdmin           = "super_admin"
	StoreManager         = "store_manager"
	ProductManager       = "product_manager"
	MarketingManager     = "marketing_manager"
	CustomerSupportLead  = "customer_support_lead"
	CustomerSupportAgent = "customer_support_agent"
	FinancialManager     = "financial_manager"
	RepairTechnician     = "repair_technician"
	UsedPhoneSpecialist  = "used_phone_specialist"
	// PartnerManager is admin staff who review and manage partners.
	PartnerManager = "partner_manager"
	// Partner is the colleague themselves — a heavily restricted backoffice
	// account that can only use the partner shop and see its own purchases.
	Partner = "partner"
)
const (
	// --- دسترسی‌های کلی پنل ---
	AccessAdminPanel = "admin_panel:access"

	// --- مدیریت داشبورد و گزارشات ---
	ViewDashboardAnalytics = "dashboard:view_analytics"
	ExportReports          = "reports:export"

	// --- مدیریت کاربران پنل (Backoffice) ---
	ListBackofficeUsers   = "backoffice_user:list"
	ViewBackofficeUser    = "backoffice_user:view"
	CreateBackofficeUser  = "backoffice_user:create"
	UpdateBackofficeUser  = "backoffice_user:update"
	DeleteBackofficeUser  = "backoffice_user:delete"
	AssignBackofficeRoles = "backoffice_user:assign_roles"
	ListPermissions       = "permission:list"
	ListRole              = "role:list"
	CreateRole            = "role:create"
	UpdateRole            = "role:update"

	// --- مدیریت مشتریان (Users) ---
	ListUsers       = "user:list"
	ViewUserDetails = "user:view_details"
	UpdateUser      = "user:update"
	BanUser         = "user:ban"

	// --- مدیریت کاتالوگ محصولات ---
	ListProducts   = "product:list"
	ViewProduct    = "product:view"
	CreateProduct  = "product:create"
	UpdateProduct  = "product:update"
	DeleteProduct  = "product:delete"
	ExportProducts = "product:export"

	// --- مدیریت موجودیت‌های وابسته به محصول ---
	ListBrands  = "brand:list"
	CreateBrand = "brand:create"
	UpdateBrand = "brand:update"
	DeleteBrand = "brand:delete"

	ListCategories = "category:list"
	CreateCategory = "category:create"
	UpdateCategory = "category:update"
	DeleteCategory = "category:delete"

	// --- مدیریت محتوا ---
	ListArticles   = "article:list"
	ViewArticle    = "article:view"
	CreateArticle  = "article:create"
	UpdateArticle  = "article:update"
	PublishArticle = "article:publish"
	DeleteArticle  = "article:delete"

	ListComment   = "comment:list"
	UpdateComment = "comment:update"

	ListReviews   = "review:list"
	ApproveReview = "review:approve"
	DeleteReview  = "review:delete"

	// --- مدیریت فروش ---
	ListOrders        = "order:list"
	ViewOrderDetails  = "order:view_details"
	UpdateOrderStatus = "order:update_status"
	CreateOrder       = "order:create"

	ListCarts       = "cart:list"
	ViewCartDetails = "cart:view_details"
	UpdateCarts     = "cart:update"
	CreateCarts     = "cart:create"

	ListDiscounts  = "discount:list"
	CreateDiscount = "discount:create"
	UpdateDiscount = "discount:update"
	DeleteDiscount = "discount:delete"

	// --- مدیریت امور مالی ---
	ListPayments       = "payment:list"
	ViewPaymentDetails = "payment:view_details"
	RefundPayment      = "payment:refund"

	// --- مدیریت خدمات ---
	ListRepairRequests = "repair:list"
	ViewRepairRequest  = "repair:view"
	UpdateRepairStatus = "repair:update_status"
	SetRepairPrice     = "repair:set_price"

	ListUsedPhoneRequests = "used_phone:list"
	ViewUsedPhoneRequest  = "used_phone:view"
	ApproveUsedPhone      = "used_phone:approve"
	RejectUsedPhone       = "used_phone:reject"

	// --- مدیریت تنظیمات ---
	UpdateSettings = "settings:update"

	// --- مدیریت همکاران (سمت کارمند) ---
	ListPartners          = "partner:list"
	ViewPartner           = "partner:view"
	ApprovePartner        = "partner:approve"
	RejectPartner         = "partner:reject"
	SuspendPartner        = "partner:suspend"
	UpdatePartnerPricing  = "partner:update_pricing"
	RecordPartnerHandover = "partner:record_handover"

	// --- پنل خرید همکار (سمت خود همکار) ---
	PartnerShopAccess = "partner_shop:access"
	PartnerShopOrder  = "partner_shop:order"

	// --- برنامه دعوت از دوستان (referral) ---
	ListReferrals  = "referral:list"
	ViewReferral   = "referral:view"
	ManageReferral = "referral:manage"

	// --- مدیریت تأمین‌کنندگان ---
	ListSuppliers  = "supplier:list"
	ViewSupplier   = "supplier:view"
	CreateSupplier = "supplier:create"
	UpdateSupplier = "supplier:update"
	DeleteSupplier = "supplier:delete"

	// --- مدیریت فاکتورهای خرید ---
	ListPurchases   = "purchase:list"
	ViewPurchase    = "purchase:view"
	CreatePurchase  = "purchase:create"
	UpdatePurchase  = "purchase:update"
	ConfirmPurchase = "purchase:confirm"
	CancelPurchase  = "purchase:cancel"

	// --- انبارداری ---
	ViewInventory   = "inventory:view"
	AdjustInventory = "inventory:adjust"
	// UpdateStockConfig covers per-variant stock mode / backorder / supplier
	// / cost settings.
	UpdateStockConfig = "inventory:update_stock_config"

	// --- هزینه‌های عملیاتی ---
	ListExpenses  = "expense:list"
	ViewExpense   = "expense:view"
	CreateExpense = "expense:create"
	UpdateExpense = "expense:update"
	DeleteExpense = "expense:delete"

	// --- صندوق و نقدینگی ---
	ViewCash   = "cash:view"
	CreateCash = "cash:create"
	DeleteCash = "cash:delete"

	// --- برگشت از فروش (مرجوعی) ---
	ListSalesReturns  = "sales_return:list"
	ViewSalesReturn   = "sales_return:view"
	CreateSalesReturn = "sales_return:create"
	RefundSalesReturn = "sales_return:refund"

	// --- گزارش‌های مالی ---
	ViewProfitReport    = "report:profit"
	ViewInventoryReport = "report:inventory"
	// ViewProfitLossReport is the full صورت سود و زیان (فروش، بهای تمام‌شده،
	// هزینه‌ها، سود خالص) — wider than the gross-profit report above.
	ViewProfitLossReport   = "report:profit_loss"
	ViewSupplierLedger     = "report:supplier_ledger"
	ViewCashReport         = "report:cash"
	ExportFinancialReports = "report:export"
)

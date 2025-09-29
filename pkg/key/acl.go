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
	ListProducts  = "product:list"
	ViewProduct   = "product:view"
	CreateProduct = "product:create"
	UpdateProduct = "product:update"
	DeleteProduct = "product:delete"

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

	ListCarts       = "cart:list"
	ViewCartDetails = "cart:view_details"

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
)

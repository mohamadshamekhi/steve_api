package key

const (
	SuperAdmin          = "Super Admin"
	StoreManager        = "Store Manager"
	MarketingSpecialist = "Marketing Specialist"
	ContentWriter       = "Content Writer"
	SeniorContentEditor = "Senior Content Editor"
	CustomerSupport     = "Customer Support"
	FinancialAccountant = "Financial Accountant"
	RepairTechnician    = "Repair Technician"
	UsedPhoneSpecialist = "Used Phone Specialist"
)
const (
	AccessAdminPanel = "admin_panel:access"

	CanListBackofficeUsers   = "backoffice_user:list"
	CanCreateBackofficeUser  = "backoffice_user:create"
	CanEditBackofficeUser    = "backoffice_user:edit"
	CanDeleteBackofficeUser  = "backoffice_user:delete"
	CanAssignBackofficeRoles = "backoffice_user:assign_roles"

	CanListUsers       = "user:list"
	CanViewUserDetails = "user:view_details"
	CanEditUser        = "user:edit"
	CanBanUser         = "user:ban"

	CanListProducts     = "product:list"
	CanCreateProduct    = "product:create"
	CanEditProduct      = "product:edit"
	CanDeleteProduct    = "product:delete"
	CanManageBrands     = "brand:manage"
	CanManageCategories = "category:manage"

	CanListArticles    = "article:list"
	CanCreateArticle   = "article:create"
	CanEditOwnArticle  = "article:edit_own"
	CanEditAllArticles = "article:edit_all"
	CanPublishArticle  = "article:publish"
	CanDeleteArticle   = "article:delete"
	CanListReviews     = "review:list"
	CanApproveReview   = "review:approve"
	CanDeleteReview    = "review:delete"

	CanViewAllOrders   = "order:view_all"
	CanEditOrderStatus = "order:edit_status"
	CanViewAllPayments = "payment:view_all"
	CanExportPayments  = "payment:export_reports"
	CanRefundPayment   = "payment:refund"
	CanManageDiscounts = "discount:manage"
	CanViewAllCarts    = "cart:view_all"

	CanViewRepairRequests    = "repair:view_requests"
	CanUpdateRepairStatus    = "repair:update_status"
	CanSetRepairPrice        = "repair:set_price"
	CanViewUsedPhoneRequests = "used_phone:view_requests"
	CanApproveUsedPhone      = "used_phone:approve"
	CanRejectUsedPhone       = "used_phone:reject"
)

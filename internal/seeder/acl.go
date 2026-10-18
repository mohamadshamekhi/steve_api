package seeder

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/key"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	"go.mongodb.org/mongo-driver/mongo"
)

var predefinedPermissions = []entity.Permission{
	{Name: key.AccessAdminPanel, Description: "اجازه ورود به پنل مدیریت"},
	{Name: key.ViewDashboardAnalytics, Description: "مشاهده تحلیل‌های داشبورد"},
	{Name: key.ExportReports, Description: "خروجی گرفتن از گزارشات"},
	{Name: key.AIAgentReport, Description: "مشاهده خلاصه تجمیعی فروش و موجودی (بدون اطلاعات مشتری و سود)"},
	{Name: key.ListBackofficeUsers, Description: "مشاهده لیست کاربران پنل"},
	{Name: key.ViewBackofficeUser, Description: "مشاهده جزئیات کاربر پنل"},
	{Name: key.CreateBackofficeUser, Description: "ایجاد کاربر جدید برای پنل"},
	{Name: key.UpdateBackofficeUser, Description: "ویرایش اطلاعات کاربران پنل"},
	{Name: key.DeleteBackofficeUser, Description: "حذف کاربر از پنل"},
	{Name: key.IssueServiceToken, Description: "صدور توکن بلندمدت برای کاربر سرویسی/ایجنت"},
	{Name: key.AssignBackofficeRoles, Description: "تخصیص و تغییر نقش‌های کاربران پنل"},
	{Name: key.ListUsers, Description: "مشاهده لیست مشتریان"},
	{Name: key.ViewUserDetails, Description: "مشاهده جزئیات پروفایل مشتری"},
	{Name: key.UpdateUser, Description: "ویرایش اطلاعات مشتریان"},
	{Name: key.BanUser, Description: "مسدود کردن حساب کاربری مشتری"},
	{Name: key.ListProducts, Description: "مشاهده لیست محصولات"},
	{Name: key.ViewProduct, Description: "مشاهده جزئیات محصول"},
	{Name: key.CreateProduct, Description: "ایجاد محصول جدید"},
	{Name: key.UpdateProduct, Description: "ویرایش محصول"},
	{Name: key.DeleteProduct, Description: "حذف محصول"},
	{Name: key.ExportProducts, Description: "خروجی اکسل از لیست محصولات"},
	{Name: key.ListBrands, Description: "مشاهده لیست برندها"},
	{Name: key.CreateBrand, Description: "ایجاد برند جدید"},
	{Name: key.UpdateBrand, Description: "ویرایش برند"},
	{Name: key.DeleteBrand, Description: "حذف برند"},
	{Name: key.ListCategories, Description: "مشاهده لیست دسته‌بندی‌ها"},
	{Name: key.CreateCategory, Description: "ایجاد دسته‌بندی جدید"},
	{Name: key.UpdateCategory, Description: "ویرایش دسته‌بندی"},
	{Name: key.DeleteCategory, Description: "حذف دسته‌بندی"},
	{Name: key.ListArticles, Description: "مشاهده لیست مقالات"},
	{Name: key.ViewArticle, Description: "مشاهده جزئیات مقاله"},
	{Name: key.CreateArticle, Description: "ایجاد مقاله جدید"},
	{Name: key.UpdateArticle, Description: "ویرایش مقاله"},
	{Name: key.PublishArticle, Description: "انتشار یا لغو انتشار مقاله"},
	{Name: key.DeleteArticle, Description: "حذف مقاله"},
	{Name: key.ListReviews, Description: "مشاهده لیست نظرات"},
	{Name: key.ApproveReview, Description: "تایید یا رد کردن نظر"},
	{Name: key.DeleteReview, Description: "حذف نظر"},
	{Name: key.ListOrders, Description: "مشاهده لیست سفارشات"},
	{Name: key.ViewOrderDetails, Description: "مشاهده جزئیات سفارش"},
	{Name: key.EditOrder, Description: "ویرایش سفارش ثبت‌شده (اقلام، مبالغ، مشتری) — فقط ادمین اصلی"},
	{Name: key.UpdateOrderStatus, Description: "تغییر وضعیت سفارش"},
	{Name: key.CreateOrder, Description: "ایجاد سفارش جدید"},
	{Name: key.ListCarts, Description: "مشاهده لیست سبدهای خرید"},
	{Name: key.ViewCartDetails, Description: "مشاهده جزئیات سبد خرید"},
	{Name: key.UpdateCarts, Description: "ویرایش سبد خرید"},
	{Name: key.CreateCarts, Description: "ایجاد سبد خرید"},
	{Name: key.ListDiscounts, Description: "مشاهده لیست تخفیف‌ها"},
	{Name: key.CreateDiscount, Description: "ایجاد تخفیف جدید"},
	{Name: key.UpdateDiscount, Description: "ویرایش تخفیف"},
	{Name: key.DeleteDiscount, Description: "حذف تخفیف"},
	{Name: key.ListPayments, Description: "مشاهده لیست پرداخت‌ها"},
	{Name: key.ViewPaymentDetails, Description: "مشاهده جزئیات پرداخت"},
	{Name: key.RefundPayment, Description: "انجام عملیات بازگشت وجه"},
	{Name: key.ListRepairRequests, Description: "مشاهده لیست درخواست‌های تعمیر"},
	{Name: key.ViewRepairRequest, Description: "مشاهده جزئیات درخواست تعمیر"},
	{Name: key.UpdateRepairStatus, Description: "تغییر وضعیت درخواست تعمیر"},
	{Name: key.SetRepairPrice, Description: "قیمت‌گذاری برای تعمیر"},
	{Name: key.ListUsedPhoneRequests, Description: "مشاهده لیست درخواست‌های فروش گوشی دست دوم"},
	{Name: key.ViewUsedPhoneRequest, Description: "مشاهده جزئیات درخواست فروش گوشی دست دوم"},
	{Name: key.ApproveUsedPhone, Description: "تایید درخواست فروش گوشی دست دوم"},
	{Name: key.RejectUsedPhone, Description: "رد کردن درخواست فروش گوشی دست دوم"},
	{Name: key.UpdateSettings, Description: "ویرایش تنظیمات کلی سایت"},
	{Name: key.UpdateComment, Description: "ویرایش کامنت"},
	{Name: key.ListComment, Description: "لیست کامنت"},
	{Name: key.ListPermissions, Description: "لیست دسترسی ها"},
	{Name: key.ListRole, Description: "لیست نقش های کاربری"},
	{Name: key.CreateRole, Description: "ایجاد نقش کاربری"},
	{Name: key.UpdateRole, Description: "ویرایش نقش کاربری"},
	{Name: key.ListPartners, Description: "مشاهده لیست همکاران و درخواست‌های همکاری"},
	{Name: key.ViewPartner, Description: "مشاهده جزئیات یک همکار"},
	{Name: key.ApprovePartner, Description: "تایید درخواست همکاری و تعیین درصد تخفیف"},
	{Name: key.RejectPartner, Description: "رد کردن درخواست همکاری"},
	{Name: key.SuspendPartner, Description: "تعلیق یا رفع تعلیق یک همکار"},
	{Name: key.UpdatePartnerPricing, Description: "ویرایش درصد تخفیف همکار و تنظیمات قیمت همکار روی محصولات"},
	{Name: key.RecordPartnerHandover, Description: "ثبت تحویل حضوری کالای سفارش همکار"},
	{Name: key.PartnerShopAccess, Description: "دسترسی همکار به پنل خرید همکاری (کاتالوگ و گزارش خرید)"},
	{Name: key.PartnerShopOrder, Description: "ثبت سفارش توسط همکار در پنل خرید همکاری"},

	{Name: key.ListSuppliers, Description: "مشاهده لیست تأمین‌کنندگان"},
	{Name: key.ViewSupplier, Description: "مشاهده جزئیات یک تأمین‌کننده"},
	{Name: key.CreateSupplier, Description: "ثبت تأمین‌کننده جدید"},
	{Name: key.UpdateSupplier, Description: "ویرایش اطلاعات تأمین‌کننده"},
	{Name: key.DeleteSupplier, Description: "حذف تأمین‌کننده"},

	{Name: key.ListPurchases, Description: "مشاهده لیست فاکتورهای خرید"},
	{Name: key.ViewPurchase, Description: "مشاهده جزئیات فاکتور خرید"},
	{Name: key.CreatePurchase, Description: "ثبت فاکتور خرید جدید"},
	{Name: key.UpdatePurchase, Description: "ویرایش فاکتور خرید پیش‌نویس"},
	{Name: key.ConfirmPurchase, Description: "تأیید فاکتور خرید و ورود کالا به انبار"},
	{Name: key.CancelPurchase, Description: "لغو فاکتور خرید و برگشت اثر آن از انبار"},

	{Name: key.ViewInventory, Description: "مشاهده موجودی و دفتر حرکت کالا"},
	{Name: key.AdjustInventory, Description: "اصلاح دستی موجودی انبار (انبارگردانی)"},
	{Name: key.UpdateStockConfig, Description: "تنظیم حالت موجودی، تأمین بیرونی و قیمت خرید محصول"},

	{Name: key.ListExpenses, Description: "مشاهده لیست هزینه‌های عملیاتی"},
	{Name: key.ViewExpense, Description: "مشاهده جزئیات یک هزینه"},
	{Name: key.CreateExpense, Description: "ثبت هزینه جدید"},
	{Name: key.UpdateExpense, Description: "ویرایش هزینه ثبت‌شده"},
	{Name: key.DeleteExpense, Description: "حذف هزینه"},

	{Name: key.ViewCash, Description: "مشاهده گردش صندوق و نقدینگی"},
	{Name: key.CreateCash, Description: "ثبت دریافت/پرداخت دستی در صندوق"},
	{Name: key.DeleteCash, Description: "حذف تراکنش دستی صندوق"},

	{Name: key.ListSalesReturns, Description: "مشاهده لیست مرجوعی‌ها"},
	{Name: key.ViewSalesReturn, Description: "مشاهده جزئیات یک مرجوعی"},
	{Name: key.CreateSalesReturn, Description: "ثبت برگشت از فروش و برگشت کالا به انبار"},
	{Name: key.RefundSalesReturn, Description: "ثبت استرداد وجه مرجوعی در صندوق"},

	{Name: key.ViewProfitReport, Description: "مشاهده گزارش سود ناخالص فروش"},
	{Name: key.ViewInventoryReport, Description: "مشاهده گزارش ارزش موجودی انبار"},
	{Name: key.ViewProfitLossReport, Description: "مشاهده صورت سود و زیان (فروش، بهای تمام‌شده، هزینه، سود خالص)"},
	{Name: key.ViewSupplierLedger, Description: "مشاهده دفتر حساب و بدهی تأمین‌کنندگان"},
	{Name: key.ViewCashReport, Description: "مشاهده گزارش دخل و خرج صندوق"},
	{Name: key.ExportFinancialReports, Description: "دریافت خروجی اکسل/CSV گزارش‌های مالی"},
	{Name: key.ListReferrals, Description: "مشاهده لیست و آمار برنامه دعوت از دوستان"},
	{Name: key.ViewReferral, Description: "مشاهده جزئیات دعوت‌های یک کاربر و کدهای تخفیف کسب‌شده"},
	{Name: key.ManageReferral, Description: "ویرایش تنظیمات برنامه دعوت (نوع و مقدار تخفیف، آستانه پاداش)"},

	{Name: key.ViewPriceReview, Description: "مشاهده لیست محصولات نیازمند بازبینی قیمت"},
	{Name: key.ManagePriceReview, Description: "تنظیم بازه بازبینی قیمت و تایید/ریست تایمر قیمت محصول"},
}

func Run(permissionRepo *mongodbpermission.DB, roleRepo *mongodbrole.DB, backooficeRepo *mongodbbackoffice.DB) error {
	ctx := context.Background()
	fmt.Println("Seeding database...")

	if err := seedPermissions(ctx, permissionRepo); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	if err := seedRoles(ctx, roleRepo); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}
	backooficeRepo.AddRoleToBackofficeUser("09011216131", key.SuperAdmin)

	fmt.Println("Database seeding completed successfully.")
	return nil
}

func seedPermissions(ctx context.Context, repo *mongodbpermission.DB) error {
	for _, p := range predefinedPermissions {
		_, err := repo.FindByName(ctx, p.Name)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			newPermission := &entity.Permission{
				Name:        p.Name,
				Description: p.Description,
			}
			newPermission.BeforeCreate()
			if err := repo.Create(ctx, newPermission); err != nil {
				return err
			}
			fmt.Printf("Permission '%s' created.\n", newPermission.Name)
		}
	}
	return nil
}

func seedRoles(ctx context.Context, repo *mongodbrole.DB) error {

	allPermissionKeys := make([]string, 0, len(predefinedPermissions))
	for _, p := range predefinedPermissions {
		allPermissionKeys = append(allPermissionKeys, p.Name)
	}

	rolesToSeed := map[string][]string{
		key.SuperAdmin: allPermissionKeys,

		key.StoreManager: {
			key.AccessAdminPanel, key.ViewDashboardAnalytics,
			key.ListUsers, key.ViewUserDetails,
			key.ListProducts, key.ViewProduct, key.ExportProducts,
			key.ListOrders, key.ViewOrderDetails, key.UpdateOrderStatus, key.CreateOrder,
			key.ListCarts, key.ViewCartDetails, key.UpdateCarts, key.CreateCarts,
			key.ListPayments, key.ViewPaymentDetails,
			key.ListDiscounts,
			key.ListSuppliers, key.ViewSupplier,
			key.ListPurchases, key.ViewPurchase, key.CreatePurchase, key.ConfirmPurchase,
			key.ViewInventory, key.AdjustInventory, key.UpdateStockConfig,
			key.ViewProfitReport, key.ViewInventoryReport, key.ViewProfitLossReport,
			key.ViewCash, key.ViewCashReport, key.CreateCash,
			key.ListExpenses, key.ViewExpense, key.CreateExpense,
			key.ListSalesReturns, key.ViewSalesReturn, key.CreateSalesReturn, key.RefundSalesReturn,
		},

		key.ProductManager: {
			key.AccessAdminPanel,
			key.ListProducts, key.ViewProduct, key.CreateProduct, key.UpdateProduct, key.DeleteProduct, key.ExportProducts,
			key.ViewInventory, key.AdjustInventory, key.UpdateStockConfig,
			key.ListSuppliers, key.ViewSupplier,
			key.ListBrands, key.CreateBrand, key.UpdateBrand, key.DeleteBrand,
			key.ListCategories, key.CreateCategory, key.UpdateCategory, key.DeleteCategory,
			key.ViewPriceReview, key.ManagePriceReview,
		},

		key.MarketingManager: {
			key.AccessAdminPanel,
			key.ListDiscounts, key.CreateDiscount, key.UpdateDiscount, key.DeleteDiscount,
			key.ListUsers, key.ViewUserDetails,
			key.ListReferrals, key.ViewReferral, key.ManageReferral,
		},

		key.CustomerSupportLead: {
			key.AccessAdminPanel,
			key.ListUsers, key.ViewUserDetails, key.UpdateUser, key.BanUser,
			key.ListOrders, key.ViewOrderDetails, key.UpdateOrderStatus, key.CreateOrder,
			key.ListCarts, key.ViewCartDetails, key.UpdateCarts, key.CreateCarts,
		},

		key.CustomerSupportAgent: {
			key.AccessAdminPanel,
			key.ListUsers, key.ViewUserDetails,
			key.ListOrders, key.ViewOrderDetails,
		},

		key.FinancialManager: {
			key.AccessAdminPanel,
			key.ListPayments, key.ViewPaymentDetails, key.RefundPayment,
			key.ExportReports,
			key.ListSuppliers, key.ViewSupplier, key.CreateSupplier, key.UpdateSupplier,
			key.ListPurchases, key.ViewPurchase, key.CreatePurchase, key.UpdatePurchase,
			key.ConfirmPurchase, key.CancelPurchase,
			key.ViewInventory, key.ViewProfitReport, key.ViewInventoryReport,
			key.ListExpenses, key.ViewExpense, key.CreateExpense, key.UpdateExpense, key.DeleteExpense,
			key.ViewCash, key.CreateCash, key.DeleteCash, key.ViewCashReport,
			key.ListSalesReturns, key.ViewSalesReturn, key.CreateSalesReturn, key.RefundSalesReturn,
			key.ViewProfitLossReport, key.ViewSupplierLedger, key.ExportFinancialReports,
		},

		key.RepairTechnician: {
			key.AccessAdminPanel,
			key.ListRepairRequests, key.ViewRepairRequest, key.UpdateRepairStatus, key.SetRepairPrice,
		},

		key.UsedPhoneSpecialist: {
			key.AccessAdminPanel,
			key.ListUsedPhoneRequests, key.ViewUsedPhoneRequest, key.ApproveUsedPhone, key.RejectUsedPhone,
		},

		key.PartnerManager: {
			key.AccessAdminPanel,
			key.ListPartners, key.ViewPartner, key.ApprovePartner, key.RejectPartner,
			key.SuspendPartner, key.UpdatePartnerPricing, key.RecordPartnerHandover,
			key.ListProducts, key.ViewProduct,
			key.ListOrders, key.ViewOrderDetails,
		},

		// The partner themselves: admin-panel login, but ONLY the partner shop.
		key.Partner: {
			key.AccessAdminPanel,
			key.PartnerShopAccess, key.PartnerShopOrder,
		},

		// key.AIAgent is the external SEO/content agent's service account.
		// It may draft (unpublished) articles and edit product SEO fields,
		// but can never delete or publish anything.
		key.AIAgent: {
			key.AccessAdminPanel,
			key.ListArticles, key.ViewArticle, key.CreateArticle, key.UpdateArticle,
			key.ListProducts, key.ViewProduct, key.UpdateProduct,
			key.AIAgentReport,
		},
	}

	for roleName, permissionKeys := range rolesToSeed {
		existingRole, err := repo.FindByName(ctx, roleName)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			newRole := &entity.Role{
				Name:           roleName,
				Description:    fmt.Sprintf("Role for %s", roleName),
				PermissionKeys: permissionKeys,
			}
			newRole.BeforeCreate()

			if err := repo.Create(ctx, newRole); err != nil {
				return err
			}
			fmt.Printf("Role '%s' created with %d permissions.\n", newRole.Name, len(permissionKeys))
			continue
		}

		// Role already exists: FindByName above only creates roles that are
		// missing entirely, so a permission added to rolesToSeed after the
		// role's first seeding would otherwise never reach the stored
		// document. Sync in any predefined permission keys the role doesn't
		// have yet, without touching permissions an admin may have added
		// manually via the role-management endpoints.
		existingSet := make(map[string]struct{}, len(existingRole.PermissionKeys))
		for _, k := range existingRole.PermissionKeys {
			existingSet[k] = struct{}{}
		}

		missing := make([]string, 0)
		for _, k := range permissionKeys {
			if _, ok := existingSet[k]; !ok {
				missing = append(missing, k)
			}
		}

		if len(missing) == 0 {
			continue
		}

		merged := append(existingRole.PermissionKeys, missing...)
		if err := repo.UpdatePermissionKeys(ctx, roleName, merged); err != nil {
			return err
		}
		fmt.Printf("Role '%s' synced with %d new permission(s): %v\n", roleName, len(missing), missing)
	}
	return nil
}

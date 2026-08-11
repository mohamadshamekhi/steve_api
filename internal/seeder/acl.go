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
	{Name: key.ListBackofficeUsers, Description: "مشاهده لیست کاربران پنل"},
	{Name: key.ViewBackofficeUser, Description: "مشاهده جزئیات کاربر پنل"},
	{Name: key.CreateBackofficeUser, Description: "ایجاد کاربر جدید برای پنل"},
	{Name: key.UpdateBackofficeUser, Description: "ویرایش اطلاعات کاربران پنل"},
	{Name: key.DeleteBackofficeUser, Description: "حذف کاربر از پنل"},
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
		},

		key.ProductManager: {
			key.AccessAdminPanel,
			key.ListProducts, key.ViewProduct, key.CreateProduct, key.UpdateProduct, key.DeleteProduct, key.ExportProducts,
			key.ListBrands, key.CreateBrand, key.UpdateBrand, key.DeleteBrand,
			key.ListCategories, key.CreateCategory, key.UpdateCategory, key.DeleteCategory,
		},

		key.MarketingManager: {
			key.AccessAdminPanel,
			key.ListDiscounts, key.CreateDiscount, key.UpdateDiscount, key.DeleteDiscount,
			key.ListUsers, key.ViewUserDetails,
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
		},

		key.RepairTechnician: {
			key.AccessAdminPanel,
			key.ListRepairRequests, key.ViewRepairRequest, key.UpdateRepairStatus, key.SetRepairPrice,
		},

		key.UsedPhoneSpecialist: {
			key.AccessAdminPanel,
			key.ListUsedPhoneRequests, key.ViewUsedPhoneRequest, key.ApproveUsedPhone, key.RejectUsedPhone,
		},
	}

	for roleName, permissionKeys := range rolesToSeed {
		_, err := repo.FindByName(ctx, roleName)
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
		}
	}
	return nil
}

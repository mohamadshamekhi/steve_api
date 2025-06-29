package seeder

import (
	"context"
	"errors"
	"fmt"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/key"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	"go.mongodb.org/mongo-driver/mongo"
)

var predefinedPermissions = []entity.Permission{
	{Name: "admin_panel:access", Description: "اجازه ورود به پنل مدیریت"},
	{Name: "backoffice_user:list", Description: "مشاهده لیست کاربران پنل مدیریت"},
	{Name: "backoffice_user:create", Description: "ایجاد کاربر جدید برای پنل مدیریت"},
	{Name: "backoffice_user:edit", Description: "ویرایش اطلاعات کاربران پنل"},
	{Name: "backoffice_user:delete", Description: "حذف کاربر از پنل مدیریت"},
	{Name: "backoffice_user:assign_roles", Description: "تخصیص و تغییر نقش‌های کاربران پنل"},
	{Name: "user:list", Description: "مشاهده لیست مشتریان"},
	{Name: "user:view_details", Description: "مشاهده جزئیات پروفایل و سفارشات یک مشتری"},
	{Name: "user:edit", Description: "ویرایش اطلاعات مشتریان"},
	{Name: "user:ban", Description: "مسدود کردن حساب کاربری یک مشتری"},
	{Name: "product:list", Description: "مشاهده لیست محصولات"},
	{Name: "product:create", Description: "ایجاد محصول جدید"},
	{Name: "product:edit", Description: "ویرایش کامل مشخصات، قیمت و موجودی محصول"},
	{Name: "product:delete", Description: "حذف محصول"},
	{Name: "brand:manage", Description: "مدیریت کامل برندها"},
	{Name: "category:manage", Description: "مدیریت کامل دسته‌بندی‌ها"},
	{Name: "article:list", Description: "مشاهده لیست مقالات وبلاگ"},
	{Name: "article:create", Description: "ایجاد مقاله جدید"},
	{Name: "article:edit_own", Description: "ویرایش مقالات پیش‌نویس خود"},
	{Name: "article:edit_all", Description: "ویرایش تمام مقالات"},
	{Name: "article:publish", Description: "انتشار یا لغو انتشار مقالات"},
	{Name: "article:delete", Description: "حذف مقاله"},
	{Name: "review:list", Description: "مشاهده لیست نظرات کاربران"},
	{Name: "review:approve", Description: "تایید یا رد کردن نظرات کاربران"},
	{Name: "review:delete", Description: "حذف نظرات تایید شده"},
	{Name: "order:view_all", Description: "مشاهده تمام سفارشات"},
	{Name: "order:edit_status", Description: "تغییر وضعیت سفارشات"},
	{Name: "payment:view_all", Description: "مشاهده تمام تراکنش‌های مالی"},
	{Name: "payment:export_reports", Description: "خروجی گرفتن از گزارشات مالی"},
	{Name: "payment:refund", Description: "انجام عملیات بازگشت وجه"},
	{Name: "discount:manage", Description: "مدیریت کامل کدهای تخفیف"},
	{Name: "cart:view_all", Description: "مشاهده سبدهای خرید فعال و رها شده"},
	{Name: "repair:view_requests", Description: "مشاهده درخواست‌های تعمیرات"},
	{Name: "repair:update_status", Description: "تغییر وضعیت یک درخواست تعمیر"},
	{Name: "repair:set_price", Description: "قیمت‌گذاری برای هزینه تعمیر"},
	{Name: "used_phone:view_requests", Description: "مشاهده درخواست‌های فروش گوشی دست دوم"},
	{Name: "used_phone:approve", Description: "تایید یک درخواست فروش"},
	{Name: "used_phone:reject", Description: "رد کردن یک درخواست فروش"},
}

func Run(permissionRepo mongodbpermission.DB, roleRepo mongodbrole.DB) error {
	ctx := context.Background()
	fmt.Println("Seeding database...")

	permissionMap, err := seedPermissions(ctx, permissionRepo)
	if err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	if err := seedRoles(ctx, roleRepo, permissionMap); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	fmt.Println("Database seeding completed successfully.")
	return nil
}

func seedPermissions(ctx context.Context, repo mongodbpermission.DB) (map[string]string, error) {
	permissionNameIdMap := make(map[string]string)

	for _, p := range predefinedPermissions {

		permission, err := repo.FindByName(ctx, p.Name)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		if err == nil {
			permissionNameIdMap[p.Name] = permission.ID
			continue
		}

		newPermission := &entity.Permission{
			Name:        p.Name,
			Description: p.Description,
		}

		newPermission.BeforeCreate()

		if err := repo.Create(ctx, newPermission); err != nil {
			return nil, err
		}

		fmt.Printf("Permission '%s' created.\n", newPermission.Name)

		permissionNameIdMap[newPermission.Name] = newPermission.ID
	}
	return permissionNameIdMap, nil
}

func seedRoles(ctx context.Context, repo mongodbrole.DB, permissionMap map[string]string) error {
	var allPermissionKeys []string
	for pKey := range permissionMap {
		allPermissionKeys = append(allPermissionKeys, pKey)
	}
	rolesToSeed := map[string][]string{
		key.SuperAdmin:          allPermissionKeys,
		key.StoreManager:        {key.AccessAdminPanel, key.CanListProducts, key.CanCreateProduct, key.CanEditProduct, key.CanDeleteProduct, key.CanManageBrands, key.CanManageCategories, key.CanViewAllOrders, key.CanEditOrderStatus, key.CanViewAllPayments, key.CanManageDiscounts, key.CanListUsers, key.CanViewUserDetails},
		key.MarketingSpecialist: {key.AccessAdminPanel, key.CanManageDiscounts, key.CanListProducts, key.CanListUsers},
		key.ContentWriter:       {key.AccessAdminPanel, key.CanListArticles, key.CanCreateArticle, key.CanEditOwnArticle},
		key.SeniorContentEditor: {key.AccessAdminPanel, key.CanListArticles, key.CanCreateArticle, key.CanEditOwnArticle, key.CanEditAllArticles, key.CanPublishArticle, key.CanDeleteArticle, key.CanListReviews, key.CanApproveReview, key.CanDeleteReview},
		key.CustomerSupport:     {key.AccessAdminPanel, key.CanViewAllOrders, key.CanEditOrderStatus, key.CanListUsers, key.CanViewUserDetails, key.CanViewAllCarts},
		key.FinancialAccountant: {key.AccessAdminPanel, key.CanViewAllPayments, key.CanExportPayments, key.CanViewAllOrders},
		key.RepairTechnician:    {key.AccessAdminPanel, key.CanViewRepairRequests, key.CanUpdateRepairStatus, key.CanSetRepairPrice},
		key.UsedPhoneSpecialist: {key.AccessAdminPanel, key.CanViewUsedPhoneRequests, key.CanApproveUsedPhone, key.CanRejectUsedPhone},
	}

	for roleName, permissionNames := range rolesToSeed {
		_, err := repo.FindByName(ctx, roleName)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			var permissionIDs []string
			for _, pName := range permissionNames {
				if id, ok := permissionMap[pName]; ok {
					permissionIDs = append(permissionIDs, id)
				}
			}

			newRole := &entity.Role{
				Name:          roleName,
				Description:   fmt.Sprintf("Role for %s", roleName),
				PermissionIDs: permissionIDs,
			}

			newRole.BeforeCreate()

			if err := repo.Create(ctx, newRole); err != nil {
				return err
			}
			fmt.Printf("Role '%s' created with %d permissions.\n", newRole.Name, len(permissionIDs))
		}
	}
	return nil
}

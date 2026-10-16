package partnerservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	productservice "github.com/abolfazlnorzad/steve/service/product"
)

// Repository is the partner persistence port.
type Repository interface {
	Create(ctx context.Context, p entity.Partner) error
	FindByID(ctx context.Context, id string) (entity.Partner, error)
	FindByUserID(ctx context.Context, userID string) (entity.Partner, error)
	FindByBackofficeUserID(ctx context.Context, backofficeUserID string) (entity.Partner, error)
	Update(ctx context.Context, p entity.Partner) error
	All(ctx context.Context, status string, page commonparam.PaginationRequest) ([]entity.Partner, int64, error)
}

// BackofficeProvisioner creates/updates the restricted backoffice account a
// partner logs in with. Satisfied by *mongodbbackoffice.DB.
type BackofficeProvisioner interface {
	GetBackOfficeUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.BackOffice, error)
	AddNewBackofficeUser(ctx context.Context, u entity.BackOffice) error
	AddRoleToBackofficeUser(phoneNumber string, role string) error
}

// PartnerOrderProvider reads the partner's own order history for the report.
type PartnerOrderProvider interface {
	GetUserOrderWithDetails(ctx context.Context, userID string, p commonparam.PaginationRequest) (orderservice.StorefrontListResult, error)
}

// ProductProvider is the slice of the product repo the partner domain needs:
// reading products for the partner catalog and writing their partner-pricing
// config from the backoffice.
type ProductProvider interface {
	FindByID(ctx context.Context, id string) (entity.Product, error)
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (productservice.ListResult, error)
	SetPartnerConfig(ctx context.Context, productID string, eligible *bool, percent *uint) error
	BulkSetPartnerEligible(ctx context.Context, productIDs []string, eligible bool) error
}

// UserProvider resolves user contact details for notifications.
type UserProvider interface {
	FindManyByIDs(ctx context.Context, ids []string) ([]entity.User, error)
}

// SMSSender sends a free-text SMS (partner approve/reject notifications).
type SMSSender interface {
	SendSimple(ctx context.Context, to string, message string) error
}

// TelegramNotifier pings the admin channel about new applications.
type TelegramNotifier interface {
	SendMessageAdmin(ctx context.Context, text string) error
}

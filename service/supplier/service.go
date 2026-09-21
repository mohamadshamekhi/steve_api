package supplierservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

// Repository is the supplier persistence port.
type Repository interface {
	Create(ctx context.Context, s entity.Supplier) error
	FindByID(ctx context.Context, id string) (entity.Supplier, error)
	FindManyByIDs(ctx context.Context, ids []string) (map[string]entity.Supplier, error)
	All(ctx context.Context, search string, isActive *bool, p commonparam.PaginationRequest) ([]entity.Supplier, int64, error)
	Update(ctx context.Context, s entity.Supplier) error
	Delete(ctx context.Context, id string) error
}

type BackofficeService struct {
	repo Repository
}

func NewBackofficeService(repo Repository) BackofficeService {
	return BackofficeService{repo: repo}
}

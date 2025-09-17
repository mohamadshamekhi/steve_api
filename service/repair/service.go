package repairservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type RepairReader interface {
	FindByID(ctx context.Context, id string) (entity.Repair, error)
	GetAll(ctx context.Context, filter cartparam.BackofficeFilterCart, pagination commonparam.PaginationRequest) (ListResult, error)
}

type AddRepair interface {
	Add(ctx context.Context, repair entity.Repair) error
}

type RepairWriter interface {
	Update(ctx context.Context, repair entity.Repair) error
	Delete(ctx context.Context, id string) error
	AddRepair
}

type RepairReadWriter interface {
	RepairReader
	RepairWriter
}

type ListResult struct {
	Repairs    []entity.Repair
	TotalCount int64
}

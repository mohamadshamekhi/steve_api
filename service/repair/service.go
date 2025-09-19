package repairservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
)

type ListResult struct {
	Repairs    []entity.Repair
	TotalCount int64
}

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error)
}

type BrandProvider interface {
	GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error)
}

type ProductProvider interface {
	GetProductsByIDs(ctx context.Context, req productparam.GetProductsByIDsRequest) (productparam.GetProductsByIDsResponse, error)
}

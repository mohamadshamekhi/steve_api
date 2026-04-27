package sellservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	userparam "github.com/abolfazlnorzad/steve/param/user"
)

type ListResult struct {
	Sells      []entity.Sell
	TotalCount int64
}

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, req userparam.GetUsersByIDsRequest) (userparam.GetUsersByIDsResponse, error)
	GetUserByID(ctx context.Context, req userparam.GetUsersRequest) (userparam.GetUsersResponse, error)
}

type BrandProvider interface {
	GetBrandsByIDs(ctx context.Context, req brandparam.GetBrandsByIDsRequest) (brandparam.GetBrandsByIDsResponse, error)
}

type ProductProvider interface {
	GetProductsByIDs(ctx context.Context, req productparam.GetProductsByIDsRequest) (productparam.GetProductsByIDsResponse, error)
}

type MelliPayamakProvider interface {
	SendPatternSMS(ctx context.Context, to string, text []string, patternKey string) error
}

type BackofficeProvider interface {
	GetSuperUserBackOffices(ctx context.Context, req backofficeparam.GetSuperUserBackofficeRequest) (backofficeparam.GetSuperUserBackofficeResponse, error)
}

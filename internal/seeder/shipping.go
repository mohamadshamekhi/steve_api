package seeder

import (
	"context"
	"errors"
	"fmt"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	mongodbshipping "github.com/abolfazlnorzad/steve/repository/mongo/shipping"
)

// FreightShippingName is the always-present shipping method used for
// oversized products that can only be delivered via freight (باربری).
// Products are flagged with entity.Product.RequiresFreightShipping; this
// seeder guarantees a matching shipping option always exists for the admin
// to price and for the storefront to offer.
const FreightShippingName = "ارسال باربری"

func SeedShipping(repo *mongodbshipping.DB) error {
	ctx := context.Background()

	_, err := repo.FindByName(ctx, FreightShippingName)
	if err == nil {
		return nil
	}
	var rErr *richerror.RichError
	if !errors.As(err, &rErr) || rErr.Kind() != richerror.KindNotFound {
		return fmt.Errorf("failed to look up freight shipping method: %w", err)
	}

	sh := entity.Shipping{
		Name:        FreightShippingName,
		Description: "ارسال با باربری برای کالاهای بزرگ که نیاز به حمل باربری دارند",
		Cost:        0,
		IsActive:    true,
	}
	sh.BeforeCreate()
	if err := repo.Create(ctx, sh); err != nil {
		return fmt.Errorf("failed to seed freight shipping method: %w", err)
	}
	fmt.Printf("Shipping method '%s' created.\n", FreightShippingName)
	return nil
}

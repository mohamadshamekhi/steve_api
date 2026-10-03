package orderservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/logger"
)

// ProductImageReader loads products (deleted ones included) to show their
// images next to order lines.
type ProductImageReader interface {
	GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error)
}

func (s *BackofficeService) SetProductImageReader(r ProductImageReader) {
	s.images = r
}

// attachThumbnails fills ProductThumbnailURL on every catalog line with one
// batched product read. Best-effort: a failure only leaves images empty.
func (s BackofficeService) attachThumbnails(ctx context.Context, orders []entity.Order) {
	if s.images == nil {
		return
	}

	seen := map[string]struct{}{}
	var ids []string
	for _, o := range orders {
		for _, it := range o.Items {
			if it.ProductID == "" {
				continue
			}
			if _, ok := seen[it.ProductID]; !ok {
				seen[it.ProductID] = struct{}{}
				ids = append(ids, it.ProductID)
			}
		}
	}
	if len(ids) == 0 {
		return
	}

	products, err := s.images.GetProductsByIDs(ctx, ids)
	if err != nil {
		logger.L().Error("failed to load product images for orders", "err", err)
		return
	}
	thumbs := make(map[string]string, len(products))
	for _, p := range products {
		thumbs[p.ID] = thumbnailOf(p)
	}

	for i := range orders {
		for j := range orders[i].Items {
			orders[i].Items[j].ProductThumbnailURL = thumbs[orders[i].Items[j].ProductID]
		}
	}
}

// thumbnailOf is the primary gallery image, else the first one.
func thumbnailOf(p entity.Product) string {
	for _, g := range p.Gallery {
		if g.IsPrimary {
			return g.URL
		}
	}
	if len(p.Gallery) > 0 {
		return p.Gallery[0].URL
	}
	return ""
}

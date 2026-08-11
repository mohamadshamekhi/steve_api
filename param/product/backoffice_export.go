package productparam

import (
	"time"

	commonparam "github.com/abolfazlnorzad/steve/param/common"
)

type ExportProductsRequest struct {
	Filter commonparam.FilterRequest
}

type ExportProductsResponse struct {
	Rows []ProductExportRow
}

type ProductExportRow struct {
	ID            string
	Name          string
	EnTitle       string
	Slug          string
	Brand         string
	Category      string
	IsPublished   bool
	IsSecondHand  bool
	MinPrice      uint64
	MaxPrice      uint64
	TotalQuantity uint64
	VariantCount  int
	SalesCount    uint
	Rating        float64
	CreatedAt     time.Time
}

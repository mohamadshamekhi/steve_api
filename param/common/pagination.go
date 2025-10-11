package commonparam

const (
	DefaultPageNumber = 1
	DefaultPageSize   = 25
)

type PaginationRequest struct {
	PageSize   uint `json:"page_size" ,query:"page_size"`
	PageNumber uint `json:"page_number" ,query:"page_number"`
}

type PaginationResponse struct {
	PageSize   uint `json:"page_size"`
	PageNumber uint `json:"page_number"`
	Total      uint `json:"total"`
}

func (p *PaginationRequest) GetPageNumber() uint {
	if p.PageNumber <= 0 {
		p.PageNumber = DefaultPageNumber
	}

	return p.PageNumber
}

func (p *PaginationRequest) GetOffset() uint {
	return (p.GetPageNumber() - 1) * p.GetPageSize()
}

func (p *PaginationRequest) GetPageSize() uint {
	validPageSizes := []uint{1, 10, 15, 25, 50, 100, 150, 200}
	for _, size := range validPageSizes {
		if p.PageSize == size {
			return size
		}
	}

	return DefaultPageSize
}

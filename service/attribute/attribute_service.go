package attributeservice

import "context"

type AttributeService struct {
}

type AttributeRepository interface {
	Create(ctx context.Context)
}

func NewAttributeService() AttributeService {
	return AttributeService{}
}

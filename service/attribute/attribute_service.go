package attributeservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type AttributeService struct {
	repo        AttributeRepository
	attrdefRepo AttributeDefinitionRepository
}
type ListResult struct {
	Attributes []entity.AttributeSet
	TotalCount int64
}
type AttributeRepository interface {
	Create(ctx context.Context, attribute entity.AttributeSet) error
	Update(ctx context.Context, attribute entity.AttributeSet) error
	All(ctx context.Context, ilter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListResult, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (entity.AttributeSet, error)
}

func NewAttributeService(repo AttributeRepository, attrdefRepo AttributeDefinitionRepository) AttributeService {
	return AttributeService{
		repo:        repo,
		attrdefRepo: attrdefRepo,
	}
}

func (s AttributeService) CreateAttribute(ctx context.Context, req attributeparam.AddAttributeRequest) (attributeparam.AddAttributeResponse, error) {
	const op = "attributeService.CreateAttribute"
	res, err := s.attrdefRepo.FindByIDs(ctx, req.AttributeDefinitionIDs)
	if err != nil {
		return attributeparam.AddAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}
	attr := entity.AttributeSet{
		Name:       req.Name,
		Attributes: res,
	}
	attr.BeforeCreate()
	if err := s.repo.Create(ctx, attr); err != nil {
		return attributeparam.AddAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return attributeparam.AddAttributeResponse{}, nil
}

func (s AttributeService) All(ctx context.Context, req attributeparam.AllAttributeRequest) (attributeparam.AllAttributeResponse, error) {
	const op = "attributeService.All"

	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return attributeparam.AllAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributeparam.AllAttributeResponse{
		Attributes: resp.Attributes,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{
			PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount),
		}}, nil
}

func (s AttributeService) Delete(ctx context.Context, req attributeparam.DeleteAttributeRequest) (attributeparam.DeleteAttributeResponse, error) {
	const op = "attributeService.Delete"
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return attributeparam.DeleteAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributeparam.DeleteAttributeResponse{}, nil
}
func (s AttributeService) Update(ctx context.Context, req attributeparam.UpdateAttributeRequest) (attributeparam.UpdateAttributeResponse, error) {
	const op = "attributeService.Update"
	attr, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return attributeparam.UpdateAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if req.Name != nil {
		attr.Name = *req.Name
	}
	if req.AttributeDefinitionIDs != nil {
		def, err := s.attrdefRepo.FindByIDs(ctx, *req.AttributeDefinitionIDs)
		if err != nil {
			return attributeparam.UpdateAttributeResponse{}, richerror.New(op).WithWarpError(err)
		}
		attr.Attributes = def
	}
	if err := s.repo.Update(ctx, attr); err != nil {
		return attributeparam.UpdateAttributeResponse{}, richerror.New(op).WithWarpError(err)
	}

	return attributeparam.UpdateAttributeResponse{}, nil
}

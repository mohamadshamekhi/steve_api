package attributeservice

import (
	"context"

	"github.com/abolfazlnorzad/steve/entity"
	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type AttributeDefinition struct {
	repo AttributeDefinitionRepository
}
type ListDefResult struct {
	Attributes []entity.AttributeDefinition
	TotalCount int64
}
type AttributeDefinitionRepository interface {
	Create(ctx context.Context, attribute entity.AttributeDefinition) error
	Update(ctx context.Context, attribute entity.AttributeDefinition) error
	All(ctx context.Context, filter commonparam.FilterRequest, page commonparam.PaginationRequest) (ListDefResult, error)
	Delete(ctx context.Context, id string) error
	FindByIDs(ctx context.Context, id []string) ([]entity.AttributeDefinition, error)
}

func NewAttributeDefinition(repo AttributeDefinitionRepository) AttributeDefinition {
	return AttributeDefinition{
		repo: repo,
	}
}

func (s AttributeDefinition) Create(ctx context.Context, req attributedefinitionparam.AddAttributeDefinitionRequest) (attributedefinitionparam.AddAttributeDefinitionResponse, error) {
	const op = "attributeDefinitionservice.Create"
	attr := entity.AttributeDefinition{
		Key:  req.Key,
		Name: req.Name,
		Unit: req.Unit,
	}
	attr.BeforeCreate()
	if err := s.repo.Create(ctx, attr); err != nil {
		return attributedefinitionparam.AddAttributeDefinitionResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributedefinitionparam.AddAttributeDefinitionResponse{}, nil
}

func (s AttributeDefinition) All(ctx context.Context, req attributedefinitionparam.AllAttributeDefinitionRequest) (attributedefinitionparam.AllAttributeDefinitionResponse, error) {
	const op = "attributeDefinitionservice.All"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return attributedefinitionparam.AllAttributeDefinitionResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributedefinitionparam.AllAttributeDefinitionResponse{
		Attributes: resp.Attributes,
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{
			PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount),
		}}, nil
}

func (s AttributeDefinition) Delete(ctx context.Context, req attributedefinitionparam.DeleteAttributeDefinitionRequest) (attributedefinitionparam.DeleteAttributeDefinitionResponse, error) {
	const op = "attributeDefinitionservice.Delete"
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return attributedefinitionparam.DeleteAttributeDefinitionResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributedefinitionparam.DeleteAttributeDefinitionResponse{}, nil
}

func (s AttributeDefinition) Update(ctx context.Context, req attributedefinitionparam.UpdateAttributeDefinitionRequest) (attributedefinitionparam.UpdateAttributeDefinitionResponse, error) {
	const op = "attributeDefinitionservice.Update"
	var ids []string
	ids = append(ids, req.ID)
	attrs, err := s.repo.FindByIDs(ctx, ids)
	if err != nil {
		return attributedefinitionparam.UpdateAttributeDefinitionResponse{}, richerror.New(op).WithWarpError(err)
	}
	attr := attrs[0]
	if req.Name != nil {
		attr.Name = *req.Name
	}
	if req.Unit != nil {
		attr.Unit = *req.Unit
	}
	if s.repo.Update(ctx, attr); err != nil {
		return attributedefinitionparam.UpdateAttributeDefinitionResponse{}, richerror.New(op).WithWarpError(err)
	}
	return attributedefinitionparam.UpdateAttributeDefinitionResponse{}, nil
}

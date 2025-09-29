package aclparam

import "github.com/abolfazlnorzad/steve/entity"

type GetAllPermissionsRequest struct {
	Search string `json:"search" query:"search"`
}

type GetAllPermissionsResponse struct {
	Permissions []entity.Permission `json:"permissions"`
}

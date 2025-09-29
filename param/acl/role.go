package aclparam

import "github.com/abolfazlnorzad/steve/entity"

type GetAllRolesRequest struct {
	Search string `json:"search" query:"search"`
}

type GetAllRolesResponse struct {
	Roles []RoleInfo `json:"roles"`
}

type RoleInfo struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Permission  []entity.Permission `json:"permission"`
}

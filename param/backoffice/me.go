package backofficeparam

import "github.com/abolfazlnorzad/steve/entity"

type MeRequest struct {
	ID string `json:"id"`
}

type MeResponse struct {
	Backoffice entity.BackOffice `json:"backoffice"`
}

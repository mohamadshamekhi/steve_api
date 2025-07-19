package attributedefinitionparam

type UpdateAttributeDefinitionRequest struct {
	ID   string  `json:"id"`
	Name *string `json:"name,omitempty"`
	Unit *string `json:"unit,omitempty"`
}

type UpdateAttributeDefinitionResponse struct {
}

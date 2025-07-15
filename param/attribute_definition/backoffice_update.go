package attributedefinitionparam

type UpdateAttributeDefinitionRequest struct {
	ID   string  `json:"id"`
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
	Unit *string `json:"unit,omitempty"`
}

type UpdateAttributeDefinitionResponse struct {
}

package attributeparam

type UpdateAttributeRequest struct {
	ID                     string    `json:"id"`
	Name                   *string   `json:"name,omitempty"`
	AttributeDefinitionIDs *[]string `json:"attribute_definition_ids,omitempty"`
}

type UpdateAttributeResponse struct {
}

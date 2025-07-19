package attributeparam

type AddAttributeRequest struct {
	Name                   string   `json:"name"`
	AttributeDefinitionIDs []string `json:"attribute_definition_ids"`
}

type AddAttributeResponse struct {
}

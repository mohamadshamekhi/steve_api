package attributedefinitionparam

type AddAttributeDefinitionRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Unit string `json:"unit,omitempty"`
}

type AddAttributeDefinitionResponse struct {
}

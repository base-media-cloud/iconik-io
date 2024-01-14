package views

// ====================================================
// iconik Objects Response Structure "GET /API/metadata/v1/views/"

// Views is the top level data structure that receives the unmarshalled payload
// response.
type Views struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ViewFields  []*ViewField `json:"view_fields"`
	Errors      interface{}  `json:"errors"`
}

// ViewField acts as a non nested struct to the ViewFields type in Views.
type ViewField struct {
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	FieldType string    `json:"field_type"`
	Options   []*Option `json:"options"`
	ReadOnly  bool      `json:"read_only"`
	Required  bool      `json:"required"`
}

// Option acts as a non nested struct to the Options type in ViewField.
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

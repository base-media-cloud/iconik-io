package metadata

type DTO struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ViewFields  []ViewFieldDTO `json:"view_fields"`
	Errors      interface{}    `json:"errors"`
}

type ViewFieldDTO struct {
	Name      string      `json:"name"`
	Label     string      `json:"label"`
	FieldType string      `json:"field_type"`
	Options   []OptionDTO `json:"options"`
	ReadOnly  bool        `json:"read_only"`
	Required  bool        `json:"required"`
}

type OptionDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

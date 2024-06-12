package metadata

// ====================================================
// iconik Objects Response Structure "GET /API/metadata/v1/views/"

// Metadata is the top level data structure that receives the unmarshalled payload
// response.
type Metadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ViewFields  []ViewField `json:"view_fields"`
	Errors      interface{} `json:"errors"`
}

// ViewField acts as a non nested struct to the ViewFields type in Metadata.
type ViewField struct {
	Name      string   `json:"name"`
	Label     string   `json:"label"`
	FieldType string   `json:"field_type"`
	Options   []Option `json:"options"`
	ReadOnly  bool     `json:"read_only"`
	Required  bool     `json:"required"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ToDTO is a method that converts a Metadata to a DTO.
func (m *Metadata) ToDTO() DTO {
	viewFieldDTOs := make([]ViewFieldDTO, len(m.ViewFields))
	for i, viewField := range m.ViewFields {
		viewFieldDTOs[i] = viewField.ToViewFieldDTO()
	}

	return DTO{
		Name:        m.Name,
		Description: m.Description,
		ViewFields:  viewFieldDTOs,
		Errors:      m.Errors,
	}
}

func (v *ViewField) ToViewFieldDTO() ViewFieldDTO {
	optionDTOs := make([]OptionDTO, len(v.Options))
	for i, option := range v.Options {
		optionDTOs[i] = option.ToOptionDTO()
	}

	return ViewFieldDTO{
		Name:      v.Name,
		Label:     v.Label,
		FieldType: v.FieldType,
		Options:   optionDTOs,
		ReadOnly:  v.ReadOnly,
		Required:  v.Required,
	}
}

func (o *Option) ToOptionDTO() OptionDTO {
	return OptionDTO{
		Label: o.Label,
		Value: o.Value,
	}
}

package views

type ViewsDTO struct {
	Name        string
	Description string
	ViewFields  []*ViewFieldDTO
	Errors      interface{}
}

type ViewFieldDTO struct {
	Name      string
	Label     string
	FieldType string
	Options   []*OptionDTO
	ReadOnly  bool
	Required  bool
}

type OptionDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

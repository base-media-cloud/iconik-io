package csv

type CSV struct {
	CSVFilesToUpdate int
	CSVMetadata      []*CSVMetadata
}

type CSVMetadata struct {
	Added                bool
	IDStruct             IDStruct
	OriginalNameStruct   OriginalNameStruct
	SizeStruct           SizeStruct
	TitleStruct          TitleStruct
	MetadataValuesStruct MetadataValuesStruct
}

type IDStruct struct {
	ID string `json:"id"`
}

type OriginalNameStruct struct {
	OriginalName string `json:"original_name"`
}

type SizeStruct struct {
	Size string `json:"size"`
}

type TitleStruct struct {
	Title string `json:"title"`
}

type MetadataValuesStruct struct {
	MetadataValues map[string]struct {
		FieldValues []FieldValue `json:"field_values"`
	} `json:"metadata_values"`
}

type FieldValue struct {
	Value string `json:"value"`
}

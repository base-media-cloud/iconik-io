package csv

type CSV struct {
	CSVMetadata []*CSVMetadata
}

type CSVMetadata struct {
	Added              bool
	IDStruct           IDStruct
	OriginalNameStruct OriginalNameStruct
	SizeStruct         SizeStruct
	TitleStruct        TitleStruct
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

package contents

// ContentsDTO is the top level data structure that receives the unmarshalled payload
// response from GET collection contents (/API/assets/v1/collections/{collection-id}/contents).
type ContentsDTO struct {
	Objects []*ObjectDTO `json:"objects"`
	Errors  interface{}  `json:"errors"`
	Pages   int
}

// ObjectDTO acts as a non nested struct to the Objects type in ContentsDTO.
type ObjectDTO struct {
	ID         string                   `json:"id"`
	Metadata   map[string][]interface{} `json:"metadata"`
	Title      string                   `json:"title"`
	Files      []*FileDTO               `json:"files"`
	ObjectType string                   `json:"object_type"`
}

// FileDTO acts as a non nested struct to the Files type in ObjectDTO.
type FileDTO struct {
	DirectoryPath string `json:"directory_path"`
	FileSetId     string `json:"file_set_id"`
	FormatId      string `json:"format_id"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	OriginalName  string `json:"original_name"`
	Size          int    `json:"size"`
	Status        string `json:"status"`
	StorageId     string `json:"storage_id"`
	StorageMethod string `json:"storage_method"`
}

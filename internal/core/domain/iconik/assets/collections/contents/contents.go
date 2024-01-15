package contents

// Contents is the top level data structure that receives the unmarshalled payload
// response from GET collection contents (/API/assets/v1/collections/{collection-id}/contents).
type Contents struct {
	Objects []*Object   `json:"objects"`
	Errors  interface{} `json:"errors"`
	Pages   int
}

// Object acts as a non nested struct to the Objects type in CollectionContents.
type Object struct {
	ID         string                   `json:"id"`
	Metadata   map[string][]interface{} `json:"metadata"`
	Title      string                   `json:"title"`
	Files      []*File                  `json:"files"`
	ObjectType string                   `json:"object_type"`
}

// File acts as a non nested struct to the Files type in Object.
type File struct {
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

func (c *Contents) ToDTO() ContentsDTO {
	var objectDTOs []*ObjectDTO

	for _, o := range c.Objects {
		objectDTO := o.ToDTO()
		objectDTOs = append(objectDTOs, &objectDTO)
	}

	return ContentsDTO{
		Objects: objectDTOs,
		Errors:  c.Errors,
		Pages:   c.Pages,
	}
}

func (o *Object) ToDTO() ObjectDTO {
	var fileDTOs []*FileDTO

	for _, f := range o.Files {
		fileDTO := f.ToDTO()
		fileDTOs = append(fileDTOs, &fileDTO)
	}

	return ObjectDTO{
		ID:         o.ID,
		Metadata:   o.Metadata,
		Title:      o.Title,
		Files:      fileDTOs,
		ObjectType: o.ObjectType,
	}
}

func (f *File) ToDTO() FileDTO {
	return FileDTO{
		DirectoryPath: f.DirectoryPath,
		FileSetId:     f.FileSetId,
		FormatId:      f.FormatId,
		Id:            f.Id,
		Name:          f.Name,
		OriginalName:  f.OriginalName,
		Size:          f.Size,
		Status:        f.Status,
		StorageId:     f.StorageId,
		StorageMethod: f.StorageMethod,
	}
}

package collections

import "time"

// Contents is the top level data structure that receives the unmarshalled payload
// response from GET collection contents (/API/assets/v1/collections/{collection-id}/contents).
type Contents struct {
	Objects []Object    `json:"objects"`
	Errors  interface{} `json:"errors"`
	Pages   int
}

// Object acts as a non nested struct to the Objects type in Contents.
type Object struct {
	ID         string                   `json:"id"`
	Metadata   map[string][]interface{} `json:"metadata"`
	Title      string                   `json:"title"`
	Files      []File                   `json:"files"`
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

// ToContentsDTO is a method that converts a Contents to a ContentsDTO.
func (c *Contents) ToContentsDTO() ContentsDTO {
	objectDTOs := make([]ObjectDTO, len(c.Objects))
	for i, object := range c.Objects {
		objectDTOs[i] = object.ToObjectDTO()
	}

	return ContentsDTO{
		Objects: objectDTOs,
		Errors:  nil,
		Pages:   0,
	}
}

// ToObjectDTO is a method that converts an Object to an ObjectDTO.
func (o *Object) ToObjectDTO() ObjectDTO {
	fileDTOs := make([]FileDTO, len(o.Files))
	for i, file := range o.Files {
		fileDTOs[i] = file.ToFileDTO()
	}

	return ObjectDTO{
		ID:         o.ID,
		Title:      o.Title,
		Files:      fileDTOs,
		ObjectType: o.ObjectType,
	}
}

// ToFileDTO is a method that converts a File to a FileDTO.
func (f *File) ToFileDTO() FileDTO {
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

// Collection is the top level data structure that receives the unmarshalled payload
// response from GET collection (/API/assets/v1/collections/{collection-id}).
type Collection struct {
	CreatedByUser     string    `json:"created_by_user"`
	CustomOrderStatus string    `json:"custom_order_status"`
	DateCreated       time.Time `json:"date_created"`
	DateModified      time.Time `json:"date_modified"`
	ID                string    `json:"id"`
	IsRoot            bool      `json:"is_root"`
	KeyframeAssetIds  []string  `json:"keyframe_asset_ids"`
	ObjectType        string    `json:"object_type"`
	Status            string    `json:"status"`
	Title             string    `json:"title"`
}

func (co *Collection) ToCollectionDTO() CollectionDTO {
	return CollectionDTO{
		CreatedByUser:     co.CreatedByUser,
		CustomOrderStatus: co.CustomOrderStatus,
		DateCreated:       co.DateCreated,
		DateModified:      co.DateModified,
		ID:                co.ID,
		IsRoot:            co.IsRoot,
		KeyframeAssetIds:  co.KeyframeAssetIds,
		ObjectType:        co.ObjectType,
		Status:            co.Status,
		Title:             co.Title,
	}
}

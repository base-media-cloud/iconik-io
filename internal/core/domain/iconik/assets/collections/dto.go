package collections

import "time"

type ContentsDTO struct {
	Objects []ObjectDTO
	Errors  interface{}
	Pages   int
}

type ObjectDTO struct {
	ID         string
	Metadata   map[string][]interface{}
	Title      string
	Files      []FileDTO
	ObjectType string
}

type FileDTO struct {
	DirectoryPath string
	FileSetId     string
	FormatId      string
	Id            string
	Name          string
	OriginalName  string
	Size          int
	Status        string
	StorageId     string
	StorageMethod string
}

type CollectionDTO struct {
	CreatedByUser     string
	CustomOrderStatus string
	DateCreated       time.Time
	DateModified      time.Time
	ID                string
	IsRoot            bool
	KeyframeAssetIds  []string
	ObjectType        string
	Status            string
	Title             string
}

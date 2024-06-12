package collections

type ContentsDTO struct {
	Objects []ObjectDTO
	Errors  interface{}
	Pages   int
}

type ObjectDTO struct {
	ID         string
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

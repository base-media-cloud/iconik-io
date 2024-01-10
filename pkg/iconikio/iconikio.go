// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

type IconikRepo interface {
	GetCol(collectionID string, pageNo int) (*Collection, error)
	GetCollection(collectionID string, pageNo int) error
	GetMetadata() error
	PrepMetadataForWriting() ([][]string, error)
	ReadCSVFile() ([][]string, error)
	WriteCSVFile(metadataFile [][]string) error
	ReadExcelFile() ([][]string, error)
	WriteExcelFile(metadataFile [][]string) error
	ProcessObjects(c *Collection, assetsMap, collectionsMap map[string]struct{}) error
	ProcessObjs(c *Collection) ([]*Object, error)
	UpdateIconik(metadataFile [][]string) error
}

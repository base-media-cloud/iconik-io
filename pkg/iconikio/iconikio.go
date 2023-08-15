// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

type IconikRepo interface {
	GetCollection(collectionID string) error
	GetMetadata() error
	PrepMetadataForWriting() ([][]string, error)
	ReadCSVFile() ([][]string, error)
	WriteCSVFile(metadataFile [][]string) error
	ReadExcelFile() ([][]string, error)
	WriteExcelFile(metadataFile [][]string) error
	ProcessObjects(c *Collection, assetsMap map[string]struct{}) error
	UpdateIconik(metadataFile [][]string) error
}

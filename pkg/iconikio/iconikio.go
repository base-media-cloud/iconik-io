// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

type IconikRepo interface {
	GetCollection(collectionID string) error
	GetMetadata() error
	WriteCSVFile() error
	ReadCSVFile() error
	ProcessObjects(c *Collection, assetsMap map[string]struct{}) error
}

// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

import (
	"encoding/csv"
)

type IconikRepo interface {
	CollectionName(collectionID string) (string, error)
	ProcessColl(collectionID string, pageNo int, w *csv.Writer) error
	WriteCollToCSV(c *Collection, w *csv.Writer) error
	Headers() [][]string
	FormatObjects(objs []*Object) ([][]string, error)
	Metadata() error
	ReadCSVFile() ([][]string, error)
	WriteCSVFile(metadataFile [][]string) error
	ReadExcelFile() ([][]string, error)
	WriteExcelFile(metadataFile [][]string) error
	UpdateIconik(metadataFile [][]string) error
}

// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

import (
	"encoding/csv"
)

type IconikRepo interface {
	GetCol(collectionID string, pageNo int, w *csv.Writer) error
	GetMetadata() error
	PrepMetadataForWriting() ([][]string, error)
	ReadCSVFile() ([][]string, error)
	WriteCSVFile(metadataFile [][]string) error
	ReadExcelFile() ([][]string, error)
	WriteExcelFile(metadataFile [][]string) error
	ProcessObjs(c *Collection, w *csv.Writer) error
	UpdateIconik(metadataFile [][]string) error
}

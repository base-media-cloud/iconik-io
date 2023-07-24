// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

type IconikRepo interface {
	GetCollection() error
	GetMetadata() error
	WriteCSVFile() error
	ReadCSVFile() error
}

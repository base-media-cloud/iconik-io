// Package iconikio contains all the of the functions required to
// pass data back and forth between the Iconik API.
package iconikio

type IconikRepo interface {
	GetCollectionAssets() error
	GetMetadata() error
	WriteCSVFile() error
	ReadCSVFile() error
	CheckAppIDAuthTokenCollectionID() error
	CheckMetadataID() error
	//CheckAssetbyID(assetID string) (int, error)
	//CheckAssetExistInCollection(assetID string) (int, error)
}

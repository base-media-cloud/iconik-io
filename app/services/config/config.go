// package config holds all configurational aspects of the service.
package config

// type Conf is the structure that holds the key variaables required
// in the execution of the service.
type Conf struct {
	IconikURL    string
	AppID        string
	AuthToken    string
	CollectionID string
	ViewID       string
	Input        string
	Output       string
	AssetIds     []string // to be appended after initial GetAsset response
	ObjectIds    []string // to be appended after initial GetAsset response
}

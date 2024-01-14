package config

// Iconik is a struct that represents the iconik config.
type Iconik struct {
	BaseURL      string
	AppID        string
	AuthToken    string
	CollectionID string
	ViewID       string
	Input        string
	Output       string
}

type Config struct {
	IconikURL string

	APIConfig        *APIConfig
	CSVFilesToUpdate int
	CSVMetadata      []*CSVMetadata
}

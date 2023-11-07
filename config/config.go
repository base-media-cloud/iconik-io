/*
Package config provides the configuration for the application.
*/
package config

type Config struct {
	AssetsPrefixURL       string
	CollectionPrefixURL   string
	MetadataViewPrefixURL string
	SearchPrefixURL       string
}

func NewConfig() Config {
	return Config{
		AssetsPrefixURL:       "/API/assets/v1/assets/",
		CollectionPrefixURL:   "/API/assets/v1/collections/",
		MetadataViewPrefixURL: "/API/metadata/v1/views/",
		SearchPrefixURL:       "/API/metadata/v1/assets/",
	}
}

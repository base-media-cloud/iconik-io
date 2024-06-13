package config

import (
	"errors"
	"flag"
	"time"
)

// Iconik is a struct that represents the iconik config.
type Iconik struct {
	BaseURL                string
	AppID                  string
	AuthToken              string
	CollectionID           string
	ViewID                 string
	OperationTimeout       time.Duration `env:"ICONIK_OPERATION_TIMEOUT,default=30s"`
	OperationRetryAttempts uint          `env:"ICONIK_OPERATION_RETRY_ATTEMPTS,default=1"`
	OperationRetryDelay    time.Duration `env:"ICONIK_OPERATION_RETRY_DELAY,default=3s"`
}

// NewIconik is a function that returns a new instance of the Iconik struct.
func NewIconik() (*Iconik, error) {
	var cfg Iconik

	flag.StringVar(&cfg.BaseURL, "iconik-url", "https://app.iconik.io", "the iconik URL")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	flag.Parse()

	if cfg.AppID == "" {
		return nil, errors.New("no App-Id provided")
	}
	if cfg.AuthToken == "" {
		return nil, errors.New("no Auth-Token provided")
	}
	if cfg.CollectionID == "" {
		return nil, errors.New("no Collection ID provided")
	}
	if cfg.ViewID == "" {
		return nil, errors.New("no Metadata View ID provided")
	}

	return &cfg, nil
}

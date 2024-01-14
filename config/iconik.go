package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

// Iconik is a struct that represents the iconik config.
type Iconik struct {
	BaseURL                string
	AppID                  string
	AuthToken              string
	CollectionID           string
	ViewID                 string
	Input                  string
	Output                 string
	OperationTimeout       time.Duration `env:"ICONIK_OPERATION_TIMEOUT,default=30s"`
	OperationRetryAttempts uint          `env:"ICONIK_OPERATION_RETRY_ATTEMPTS,default=1"`
	OperationRetryDelay    time.Duration `env:"ICONIK_OPERATION_RETRY_DELAY,default=3s"`
	Limit                  int           `env:"ICONIK_OBJECT_LIMIT,default=1000"`
}

// NewIconik is a function that returns a new instance of the Iconik struct.
func NewIconik() (*Iconik, error) {
	var cfg Iconik

	flag.StringVar(&cfg.BaseURL, "iconik-url", "app.iconik.io", "the iconik URL")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	flag.StringVar(&cfg.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cfg.Output, "output", "", "Output mode - requires path to save CSV file")
	ver := flag.Bool("version", false, "Print version")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *ver {
		return nil, errors.New("version selected")
	}
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
	if cfg.Input == "" && cfg.Output == "" {
		versionInfo()
		return nil, errors.New("neither input or output mode selected")
	}

	return &cfg, nil
}

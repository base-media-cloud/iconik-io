package config

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
)

var version string

// App is a struct that represents the app config.
type App struct {
	Type                   string
	Input                  string
	Output                 string
	BaseURL                string
	AppID                  string
	AuthToken              string
	CollectionID           string
	ViewID                 string
	OperationTimeout       time.Duration `env:"OPERATION_TIMEOUT,default=30s"`
	OperationRetryAttempts uint          `env:"OPERATION_RETRY_ATTEMPTS,default=1"`
	OperationRetryDelay    time.Duration `env:"OPERATION_RETRY_DELAY,default=3s"`
	PerPage                int           `env:"PER_PAGE,default=150"`
	Version                string        `env:"VERSION"`
	Build                  string        `env:"BUILD"`
	Year                   int           `env:"YEAR,default=2024"`
}

// NewApp is a function that returns a new instance of the App struct.
func NewApp() (*App, error) {
	fmt.Println(version)
	var cfg App
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cfg.Output, "output", "", "Output mode - requires path to save CSV file")
	flag.StringVar(&cfg.BaseURL, "iconik-url", "https://app.iconik.io", "the iconik URL")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	ver := flag.Bool("version", false, "Print version")
	flag.Parse()

	cfg.Version = version

	if *ver {
		cfg.Print()
	}

	if flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if cfg.Input != "" && cfg.Output != "" {
		fmt.Println("both input or output mode selected. Please only select one.")
		return nil, nil
	}

	if cfg.Input == "" && cfg.Output == "" {
		fmt.Println("neither input or output mode selected")
		return nil, nil
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

	if cfg.Input != "" {
		cfg.Type = "input"
	}

	if cfg.Output != "" {
		cfg.Type = "output"
	}

	return &cfg, nil
}

// Print prints the version info.
func (a *App) Print() {
	fmt.Printf(`
base iconik-io
iconik CSV read/write tool
Version: %s | Build: %s
Copyright © %d Base Media Cloud Limited
https://base-mc.com
`, a.Version, a.Build, a.Year)
	os.Exit(1)
}

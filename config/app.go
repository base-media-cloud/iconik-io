package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

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
	OperationTimeout       time.Duration
	OperationRetryAttempts uint
	OperationRetryDelay    time.Duration
}

// NewApp is a function that returns a new instance of the App struct.
func NewApp(build, version string) (*App, error) {
	var cfg App

	cfg.OperationTimeout = time.Second * 30
	cfg.OperationRetryAttempts = 1
	cfg.OperationRetryDelay = time.Second * 3

	flag.StringVar(&cfg.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cfg.Output, "output", "", "Output mode - requires path to save CSV file")
	flag.StringVar(&cfg.BaseURL, "iconik-url", "https://app.iconik.io", "the iconik URL")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	ver := flag.Bool("version", false, "Print version")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *ver {
		versionInfo(build, version)
		return nil, nil
	}

	if cfg.Input != "" && cfg.Output != "" {
		fmt.Println("both input or output mode selected. Please only select one.")
		versionInfo(build, version)
		return nil, nil
	}

	if cfg.Input == "" && cfg.Output == "" {
		fmt.Println("neither input or output mode selected")
		versionInfo(build, version)
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

func versionInfo(build, version string) {
	fmt.Printf(`
base iconik-io
iconik CSV read/write tool
Version: %s | Build: %s
Copyright © 2023 Base Media Cloud Limited
https://base-mc.com
`, version, build)
}

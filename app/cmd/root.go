package cmd

import (
	"flag"

	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/iconikio"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	app   Application
	build = "develop"
)

type Application struct {
	Logger *zap.SugaredLogger
	Iconik iconikio.IconikRepo
}

func Execute(l *zap.SugaredLogger) error {

	// Add logger to part of our Application struct and log
	app.Logger = l
	app.Logger.Infow("starting service", zapcore.Field{
		Key:    "build",
		Type:   zapcore.StringType,
		String: build,
	})
	defer app.Logger.Infow("shutdown complete")

	// Parse command line flags, store in Config struct
	cfg, err := argParse()
	if err != nil {
		return err
	}

	// Create new Iconik Client struct
	iconikClient := iconikio.New(cfg)

	// Populate Iconik URL struct
	iconikClient.NewAPIConfig()

	// Attach Iconik Client to Iconik Repo interface
	app.Iconik = &iconikio.Iconik{IconikClient: iconikClient}

	// Validate App ID, Auth Token and Collection ID
	err = app.Iconik.CheckAppIDAuthTokenCollectionID()
	if err != nil {
		return err
	}

	// Validate Metadata ID
	err = app.Iconik.CheckMetadataID()
	if err != nil {
		return err
	}

	if cfg.Output != "" {
		// User has chosen CSV output:
		// Get Assets
		err = app.Iconik.GetCollectionAssets()
		if err != nil {
			return err
		}

		// Build CSV and output
		err = app.Iconik.WriteCSVFile()
		if err != nil {
			return err
		}
	}

	if cfg.Input != "" {
		// User has chosen CSV input:
		// Read CSV file and update metadata and title on Iconik API
		err := app.Iconik.ReadCSVFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func argParse() (*iconikio.Config, error) {

	var cfg iconikio.Config

	flag.StringVar(&cfg.IconikURL, "iconik-url", "https://preview.iconik.cloud", "the iconik URL (default https://preview.iconik.cloud)")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	flag.StringVar(&cfg.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cfg.Output, "output", "", "Output mode - requires path to save CSV file")
	flag.Parse()

	if cfg.AppID == "" {
		app.Logger.Fatalw("No App-Id provided")
	}
	if cfg.AuthToken == "" {
		app.Logger.Fatalw("No Auth-Token provided")
	}
	if cfg.CollectionID == "" {
		app.Logger.Fatalw("No Collection ID provided")
	}
	if cfg.ViewID == "" {
		app.Logger.Fatalw("No Metadata View ID provided")
	}
	if cfg.Input == "" && cfg.Output == "" {
		app.Logger.Fatalw("Neither input or output mode selected. Please select one.")
	}

	return &cfg, nil
}

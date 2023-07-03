package cmd

import (
	"flag"
	"log"

	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/iconikio"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

type Application struct {
	log    *zap.SugaredLogger
	Iconik iconikio.IconikRepo
}

func Execute() error {

	app := Application{}

	log, err := logger.New("PD-ICONIK-IO-RD")
	if err != nil {
		return err
	}
	defer log.Sync()

	log.Infow("starting service", zapcore.Field{
		Key:    "build",
		Type:   zapcore.StringType,
		String: build,
	})
	defer log.Infow("shutdown complete")
	app.log = log

	cfg, err := argParse()
	if err != nil {
		return err
	}

	iconikClient := iconikio.New(cfg)

	app.Iconik = &iconikio.Iconik{IconikClient: iconikClient}

	err = app.Iconik.CheckAppIDAuthTokenCollectionID()
	if err != nil {
		return err
	}

	err = app.Iconik.CheckMetadataID()
	if err != nil {
		return err
	}

	if cfg.Output != "" {
		// User has chosen CSV output

		// Get Assets
		err = app.Iconik.GetCollectionAssets()
		if err != nil {
			return err
		}

		// Get CSV Headers
		columnsName, columnsLabel, err := app.Iconik.GetCSVColumnsFromView()
		if err != nil {
			return err
		}

		// Build CSV and output
		err = app.Iconik.BuildCSVFile(columnsName, columnsLabel)
		if err != nil {
			return err
		}
	}

	if cfg.Input != "" {
		// User has chosen CSV input
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
		log.Fatal("No App-Id provided")
	}
	if cfg.AuthToken == "" {
		log.Fatal("No Auth-Token provided")
	}
	if cfg.CollectionID == "" {
		log.Fatal("No Collection ID provided")
	}
	if cfg.ViewID == "" {
		log.Fatal("No Metadata View ID provided")
	}
	if cfg.Input == "" && cfg.Output == "" {
		log.Fatal("Neither input or output mode selected. Please select one.")
	}

	return &cfg, nil
}

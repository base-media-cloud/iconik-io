package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/iconikio"
)

var (
	app     Application
	build   string
	version = "0.03"
)

type Application struct {
	Logger *zap.SugaredLogger
	Iconik iconikio.IconikRepo
}

func Execute(l *zap.SugaredLogger, appCfg config.Config) error {
	app.Logger = l
	app.Logger.Infow("starting service", zapcore.Field{
		Key:    "build",
		Type:   zapcore.StringType,
		String: build,
	})
	defer app.Logger.Infow("shutdown complete")

	cfg, err := argParse()
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}

	// Create new Iconik Client struct
	iconikClient := iconikio.New(cfg)

	// Populate Iconik URL struct
	iconikClient.NewAPIConfig(appCfg)

	// Attach Iconik Client to Iconik Repo interface
	app.Iconik = &iconikio.Iconik{IconikClient: iconikClient}

	// Get Metadata using given Metadata ID
	err = app.Iconik.GetMetadata()
	if err != nil {
		return err
	}

	// Get Collection using given Collection ID
	err = app.Iconik.GetCollection(cfg.CollectionID, 1)
	if err != nil {
		return err
	}

	assetsMap := make(map[string]struct{})
	collectionsMap := make(map[string]struct{})
	err = app.Iconik.ProcessObjects(iconikClient.Collection, assetsMap, collectionsMap)
	if err != nil {
		return err
	}

	if cfg.Output != "" && cfg.Excel {
		// Build CSV and output
		metadataFile, err := app.Iconik.PrepMetadataForWriting()
		if err != nil {
			return err
		}

		err = app.Iconik.WriteExcelFile(metadataFile)
		if err != nil {
			return err
		}
	}

	if cfg.Output != "" && cfg.CSV {
		// Build CSV and output
		metadataFile, err := app.Iconik.PrepMetadataForWriting()
		if err != nil {
			return err
		}

		err = app.Iconik.WriteCSVFile(metadataFile)
		if err != nil {
			return err
		}
	}

	if filepath.Ext(cfg.Input) == ".xlsx" {
		fmt.Println()
		fmt.Println("Inputting data from provided Excel file:")

		excelData, err := app.Iconik.ReadExcelFile()
		if err != nil {
			return err
		}
		err = app.Iconik.UpdateIconik(excelData)
		if err != nil {
			return err
		}
	}

	if filepath.Ext(cfg.Input) == ".csv" {
		fmt.Println()
		fmt.Println("Inputting data from provided CSV file:")

		csvData, err := app.Iconik.ReadCSVFile()
		if err != nil {
			return err
		}
		err = app.Iconik.UpdateIconik(csvData)
		if err != nil {
			return err
		}
	}

	return nil
}

func argParse() (*iconikio.Config, error) {

	var cfg iconikio.Config

	flag.StringVar(&cfg.IconikURL, "iconik-url", "app.iconik.io", "the iconik URL")
	flag.StringVar(&cfg.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cfg.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cfg.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cfg.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	flag.StringVar(&cfg.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cfg.Output, "output", "", "Output mode - requires path to save CSV file")
	flag.BoolVar(&cfg.Excel, "excel", false, "Select Excel output")
	flag.BoolVar(&cfg.CSV, "csv", false, "Select CSV output")
	ver := flag.Bool("version", false, "Print version")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *ver {
		versionInfo()
		return nil, nil
	}

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
		app.Logger.Infoln("Neither input or output mode selected")
		versionInfo()
		return nil, nil
	}
	if cfg.Output != "" && !cfg.Excel && !cfg.CSV {
		app.Logger.Infoln("Neither excel or csv file format selected")
		versionInfo()
		return nil, nil
	}

	return &cfg, nil
}

func versionInfo() {
	fmt.Printf(`
base iconik-io
iconik CSV read/write tool
Version: %s | Build: %s
Copyright © 2023 Base Media Cloud Limited
https://base-mc.com
`, version, build)
}

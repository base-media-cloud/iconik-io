/*
Package cmd executes the commands required to run the application.
*/
package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/iconikio"
)

var (
	app     Application
	build   string
	version = "0.05"
)

type Application struct {
	Logger zerolog.Logger
	Iconik iconikio.IconikRepo
}

func Execute(l zerolog.Logger) error {
	app.Logger = l
	app.Logger.Info().Str("build", build).Msg("starting service")

	iconikCfg, err := config.NewIconik()
	if err != nil {
		return handleErrorResponse(err)
	}

	ctx := app.Logger.WithContext(context.Background())

	req := api.New(&http.Client{})
	iconikAPI := iconik.New(iconikCfg, req)

	views, err := iconikAPI.GetMetadataViews(ctx, iconik.MetadataViewsPath, iconikCfg.ViewID)
	if err != nil {
		return err
	}

	if iconikCfg.Input != "" {
		fmt.Println("\nInputting data from provided CSV file:")

		csvData, err := app.Iconik.ReadCSVFile()
		if err != nil {
			return err
		}
		err = app.Iconik.UpdateIconik(csvData)
		if err != nil {
			return err
		}
	}

	if iconikCfg.Output != "" {
		collectionName, err := app.Iconik.CollectionName(iconikCfg.CollectionID)
		if err != nil {
			return err
		}

		today := time.Now().Format("2006-01-02_150405")
		filename := fmt.Sprintf("%s_%s_Report_%s.csv", iconikCfg.CollectionID, collectionName, today)
		filePath := iconikCfg.Output + filename

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		w := csv.NewWriter(f)
		if err = w.WriteAll(app.Iconik.Headers()); err != nil {
			return err
		}
		if err = app.Iconik.ProcessColl(iconikCfg.CollectionID, 1, w); err != nil {
			return err
		}
	}

	return nil
}

func handleErrorResponse(err error) error {
	switch {
	case errors.Is(err, errors.New("neither input or output mode selected")), errors.Is(err, errors.New("version selected")):
		versionInfo()
		return nil
	default:
		return err
	}
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

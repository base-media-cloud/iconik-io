/*
Package cmd executes the commands required to run the application.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/iconikio"
	"go.uber.org/zap"
)

var (
	app     Application
	Build   string
	Version = "0.05b"
)

type Application struct {
	Logger *zap.SugaredLogger
	Iconik iconikio.IconikRepo
}

func Execute(appCfg config.Config, l zerolog.Logger) error {

	iconikCfg, err := config.NewIconik()
	if err != nil {
		l.Fatal().Err(err).Msg("error creating iconik config")
	}

	req := api.New(&http.Client{})

	// Create new Iconik Client struct
	iconikClient := iconikio.New(cfg)

	// Populate Iconik URL struct
	iconikClient.NewAPIConfig(appCfg)

	// Attach Iconik Client to Iconik Repo interface
	app.Iconik = &iconikio.Iconik{IconikClient: iconikClient}

	// Get Metadata using given Metadata ID
	err = app.Iconik.Metadata()
	if err != nil {
		return err
	}

	if iconikCfg.Output != "" {
		collectionName, err := app.Iconik.CollectionName(iconikCfg.CollectionID)
		if err != nil {
			return err
		}

		today := time.Now().Format("2006-01-02_150405")
		filename := fmt.Sprintf("%s_%s_Report_%s.csv", iconikCfg.CollectionID, collectionName, today)
		filePath := iconikClient.Config.Output + filename

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

	if iconikCfg.Input != "" {
		err = app.Iconik.GetCollection(iconikCfg.CollectionID, 1)
		if err != nil {
			return err
		}

		assetsMap := make(map[string]struct{})
		collectionsMap := make(map[string]struct{})
		err = app.Iconik.ProcessObjects(iconikClient.Collection, assetsMap, collectionsMap)
		if err != nil {
			return err
		}

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

	return nil
}

package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	inputsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/input"
	"github.com/rs/zerolog"
)

// AppType is the app type which determines if the app should run in input mode.
const AppType = "input"

// Run runs the functions to input data from a csv into iconik.
func Run(cfg *config.App, inputSvc *inputsvc.Svc, l zerolog.Logger) error {
	fmt.Println("\nInputting data from provided CSV file...")

	ctx := l.WithContext(context.Background())

	view, err := inputSvc.GetMetadataView(ctx, cfg.ViewID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to retrieve metadata view")
		return err
	}

	csvData, err := inputSvc.ReadCSVFile(cfg)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to read csv file")
		return err
	}

	csvHeaders := csvData[0]
	if csvHeaders[0] != "id" || csvHeaders[1] != "original_name" || csvHeaders[2] != "size" || csvHeaders[3] != "title" {
		fmt.Println(csvHeaders)
		return errors.New("CSV file not properly formatted for Iconik")
	}

	matchingData, nonMatchingHeaders, err := inputSvc.MatchCSVtoView(view.ViewFields, csvData)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		fmt.Printf(`
Some columns from the file provided have not been included in the upload to Iconik, 
as they are not part of the metadata view provided. 

Please see below for the headers of the columns not included:
`)
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	csvFilesToUpdate := len(matchingData) - 2
	fmt.Println("Amount of files to update:", csvFilesToUpdate)

	notAdded, err := inputSvc.ProcessAssets(ctx, matchingData, cfg.CollectionID, cfg.ViewID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to write csv to iconik")
		return err
	}

	fmt.Printf("Assets successfully updated: %d of %d\n", csvFilesToUpdate-len(notAdded), csvFilesToUpdate)
	if len(notAdded) > 0 {
		fmt.Println("Some assets failed to update:")
		for assetID := range notAdded {
			fmt.Printf("Asset ID: %s\n", assetID)
		}
	}

	return nil
}

package input

import (
	"context"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	collDomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	inputsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/input"
	"github.com/rs/zerolog"
)

// AppType is the app type which determines if the app should run in input mode.
const AppType = "input"

// Run runs the functions to input data from a csv into iconik.
func Run(cfg *config.App, inputSvc *inputsvc.Svc, l zerolog.Logger) error {
	ctx := l.WithContext(context.Background())

	var err error

	view, err := inputSvc.GetMetadataView(ctx, cfg.ViewID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to retrieve metadata view")
		return err
	}

	var objects []collDomain.ObjectDTO
	objects, err = inputSvc.GetCollectionObjects(ctx, cfg.CollectionID, 1, objects)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to retrieve collection contents")
		return err
	}

	assetsMap := make(map[string]struct{})
	collectionsMap := make(map[string]struct{})
	var assets []collDomain.ObjectDTO
	assets, err = inputSvc.ProcessObjects(ctx, assets, objects, assetsMap, collectionsMap)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to process objects")
		return err
	}

	fmt.Println("\nInputting data from provided CSV file:")

	csvData, err := inputSvc.ReadCSVFile(cfg)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to read csv file")
		return err
	}
	err = inputSvc.UpdateIconik(ctx, view.ViewFields, assets, csvData, cfg)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to update iconik")
		return err
	}

	return nil
}

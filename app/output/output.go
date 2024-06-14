package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	outputsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/output"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// AppType is the app type which determines if the app should run in output mode.
const AppType = "output"

// Run runs the functions to output data from iconik to a csv.
func Run(cfg *config.App, outputSvc *outputsvc.Svc, l zerolog.Logger) error {
	ctx := l.WithContext(context.Background())
	fmt.Println("Running output...")

	var err error

	view, err := outputSvc.GetMetadataView(ctx, cfg.ViewID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to retrieve metadata view")
		return err
	}

	coll, err := outputSvc.GetCollection(ctx, cfg.CollectionID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to retrieve collection")
		return err
	}

	filePath := cfg.Output + fmt.Sprintf("%s_%s_Report_%s.csv", cfg.CollectionID, coll.Title, time.Now().Format("2006-01-02_150405"))

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err = w.WriteAll(outputSvc.Headers(view.ViewFields)); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to write headers to csv")
		return err
	}

	if err = outputSvc.ProcessPage(ctx, view.ViewFields, cfg.CollectionID, []interface{}{}, w); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to write assets to csv")
		return err
	}

	fmt.Println("Output complete. CSV created at " + filePath)

	return nil
}

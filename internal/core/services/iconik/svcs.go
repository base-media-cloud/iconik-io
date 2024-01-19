package iconik

import (
	"context"
	"encoding/csv"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"github.com/rs/zerolog"
)

// Svc is a struct that implements the iconik servicer ports.
type Svc struct {
	metaSvc metadata.Servicer
	collSvc collections.Servicer
	csvSvc  csvsvc.Servicer
}

// New is a function that returns a new instance of iconik Svc struct.
func New(
	metaSvc metadata.Servicer,
	collSvc collections.Servicer,
) *Svc {
	return &Svc{
		metaSvc: metaSvc,
		collSvc: collSvc,
	}
}

// ProcessCollection hits the iconik endpoint, gets system domains and creates lines in the database.
func (svc *Svc) ProcessCollection(ctx context.Context, path, collectionID string, pageNo int, w *csv.Writer) error {
	objDTOs, err := svc.collSvc.GetCollContents(ctx, path, collectionID, pageNo)
	if err != nil {
		zerolog.Ctx(ctx).Info().
			Err(err).
			Msg("error getting collection contents")
		return err
	}

	if err = svc.csvSvc.WriteObjsToCSV(objDTOs, w); err != nil {
		zerolog.Ctx(ctx).Info().
			Err(err).
			Msg("error writing objects to csv")
		return err
	}

	return nil
}

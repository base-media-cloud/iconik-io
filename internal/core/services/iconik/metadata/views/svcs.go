package views

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata/views"
	"github.com/rs/zerolog"
)

// API is an interface that defines the operations that can be performed on the billing endpoint.
type API interface {
	GetMetadataViews(ctx context.Context, path string) ([]views.ViewFieldDTO, error)
}

// Svc is a struct that implements the metadataports.Servicer interface.
type Svc struct {
	api API
	cfg *config.Iconik
}

// New is a function that returns a new instance of the Svc struct.
func New(
	a API,
	cfg *config.Iconik,
) *Svc {
	return &Svc{
		api: a,
		cfg: cfg,
	}
}

func (s *Svc) GetMetadataViews(ctx context.Context, path string) ([]views.ViewFieldDTO, error) {
	dtos, err := s.api.GetMetadataViews(ctx, path)
	if err != nil {
		return nil, err
	}

	if len(dtos) == 0 {
		zerolog.Ctx(ctx).Info().
			Str("path", iconik.MetadataViewsPath).
			Msg("no metadata views returned from iconik api")
		return nil, nil
	}

	return dtos, nil
}

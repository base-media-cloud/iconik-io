package assets

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
)

// API is an interface that defines the operations that can be performed on the metadata assets endpoint.
type API interface {
	UpdateMetadataInAsset(ctx context.Context, path, viewID, assetID string, payload []byte) (metadata.DTO, error)
}

type Svc struct {
	api API
}

// New is a function that returns a new instance of the Svc struct.
func New(
	api API,
) *Svc {
	return &Svc{
		api: api,
	}
}

// UpdateMetadataInAsset updates an assets metadata in the iconik api.
func (s *Svc) UpdateMetadataInAsset(ctx context.Context, path, viewID, assetID string, payload []byte) (metadata.DTO, error) {
	dto, err := s.api.UpdateMetadataInAsset(ctx, path, assetID)
	if err != nil {
		return assets.DTO{}, err
	}

	return dto, nil
}

// UpdateAsset updates an asset in the iconik api.
func (s *Svc) UpdateAsset(ctx context.Context, path, assetID string, payload []byte) (assets.DTO, error) {
	dto, err := s.api.PatchAsset(ctx, path, assetID, payload)
	if err != nil {
		return assets.DTO{}, err
	}

	return dto, nil
}

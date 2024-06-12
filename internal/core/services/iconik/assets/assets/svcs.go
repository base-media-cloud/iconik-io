package assets

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/assets"
)

// API is an interface that defines the operations that can be performed on the assets endpoint.
type API interface {
	GetAsset(ctx context.Context, path, assetID string) (assets.DTO, error)
	PatchAsset(ctx context.Context, path, assetID string, payload []byte) (assets.DTO, error)
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

// GetAsset gets an asset from the iconik api.
func (s *Svc) GetAsset(ctx context.Context, path, assetID string) (assets.DTO, error) {
	dto, err := s.api.GetAsset(ctx, path, assetID)
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

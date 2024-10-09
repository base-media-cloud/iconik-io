package assets

import (
	"context"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/assets"
	"github.com/google/uuid"
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

// ValidateAsset validates an asset in the iconik api.
func (s *Svc) ValidateAsset(ctx context.Context, assetID string) error {
	_, err := uuid.Parse(assetID)
	if err != nil {
		return errors.New("not a valid asset ID")
	}

	_, err = s.api.GetAsset(ctx, iconik.AssetsPath, assetID)
	if err != nil {
		return err
	}

	return fmt.Errorf("asset %s does not exist in given collection id", assetID)
}

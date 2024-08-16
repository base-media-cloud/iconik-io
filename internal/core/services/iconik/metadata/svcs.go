package metadata

import (
	"context"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
)

// API is an interface that defines the operations that can be performed on the metadata endpoint.
type API interface {
	UpdateMetadataInAsset(ctx context.Context, path, viewID, assetID string, payload []byte) (metadata.DTO, error)
	GetMetadataView(ctx context.Context, path, viewID string) (metadata.DTO, error)
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
	dto, err := s.api.UpdateMetadataInAsset(ctx, path, viewID, assetID, payload)
	if err != nil {
		return metadata.DTO{}, err
	}

	return dto, nil
}

// GetMetadataView gets a metadata view from the iconik api.
func (s *Svc) GetMetadataView(ctx context.Context, path, viewID string) (metadata.DTO, error) {
	dto, err := s.api.GetMetadataView(ctx, path, viewID)
	if err != nil {
		return metadata.DTO{}, err
	}

	if dto.Errors != nil {
		return metadata.DTO{}, fmt.Errorf("%v", dto.Errors)
	}

	return dto, nil
}

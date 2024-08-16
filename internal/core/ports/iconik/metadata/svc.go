package metadata

//go:generate mockgen -source svc.go -destination=../../../../../mocks/metadata_mocks/svc.go -package=metadata_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	UpdateMetadataInAsset(ctx context.Context, path, viewID, assetID string, payload []byte) (metadata.DTO, error)
	GetMetadataView(ctx context.Context, path, viewID string) (metadata.DTO, error)
}

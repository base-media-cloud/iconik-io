package assets

//go:generate mockgen -source svc.go -destination=../../../../../mocks/assets_mocks/svc.go -package=assets_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/assets"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	GetAsset(ctx context.Context, path, assetID string) (assets.DTO, error)
	UpdateAsset(ctx context.Context, path, assetID string, payload []byte) (assets.DTO, error)
	ValidateAsset(ctx context.Context, assetID string) error
}

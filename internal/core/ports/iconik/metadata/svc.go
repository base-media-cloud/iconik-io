package metadata

//go:generate mockgen -source svc.go -destination=../../../../../mocks/metadata_mocks/svc.go -package=metadata_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata/views"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	GetMetadataViews(ctx context.Context, path, viewID string) ([]views.ViewFieldDTO, error)
}

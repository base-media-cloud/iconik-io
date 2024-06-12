package collections

//go:generate mockgen -source svc.go -destination=../../../../../mocks/collections_mocks/svc.go -package=collections_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	GetCollectionContents(ctx context.Context, path, collectionID string) (collections.ContentsDTO, error)
}

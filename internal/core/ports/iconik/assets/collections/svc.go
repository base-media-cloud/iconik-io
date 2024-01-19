package collections

//go:generate mockgen -source svc.go -destination=../../../../../mocks/collections_mocks/svc.go -package=collections_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections/contents"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	GetCollContents(ctx context.Context, path, collectionID string, pageNo int) ([]contents.ObjectDTO, error)
	GetCollection(ctx context.Context, path, collectionID string) (collections.DTO, error)
}

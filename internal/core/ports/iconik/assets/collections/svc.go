package collections

//go:generate mockgen -source svc.go -destination=../../../../../mocks/collections_mocks/svc.go -package=collections_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	GetContents(ctx context.Context, path, collectionID string, pageNo int) (collections.ContentsDTO, error)
	GetCollection(ctx context.Context, path, collectionID string) (collections.CollectionDTO, error)
}

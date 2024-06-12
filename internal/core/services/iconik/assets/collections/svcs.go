package collections

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
)

// API is an interface that defines the operations that can be performed on the collections endpoint.
type API interface {
	GetCollectionContents(ctx context.Context, path, collectionID string) (collections.ContentsDTO, error)
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

// GetContents gets the collection contents from the iconik api.
func (s *Svc) GetContents(ctx context.Context, path, collectionID string) (collections.ContentsDTO, error) {
	dtos, err := s.api.GetCollectionContents(ctx, path, collectionID)
	if err != nil {
		return collections.ContentsDTO{}, err
	}

	return dtos, nil
}

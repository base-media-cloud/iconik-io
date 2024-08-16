package collections

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	"strconv"
)

// API is an interface that defines the operations that can be performed on the collections endpoint.
type API interface {
	GetCollectionContents(ctx context.Context, path, collectionID string, queryParams map[string]string) (collections.ContentsDTO, error)
	GetCollection(ctx context.Context, path, collectionID string) (collections.CollectionDTO, error)
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
func (s *Svc) GetContents(ctx context.Context, path, collectionID string, pageNo int) (collections.ContentsDTO, error) {
	queryParams := map[string]string{
		"per_page": "500",
		"page":     strconv.Itoa(pageNo),
	}

	dtos, err := s.api.GetCollectionContents(ctx, path, collectionID, queryParams)
	if err != nil {
		return collections.ContentsDTO{}, err
	}

	return dtos, nil
}

// GetCollection takes a collection ID and returns the collection.
func (s *Svc) GetCollection(ctx context.Context, path, collectionID string) (collections.CollectionDTO, error) {
	dto, err := s.api.GetCollection(ctx, path, collectionID)
	if err != nil {
		return collections.CollectionDTO{}, err
	}

	return dto, nil
}

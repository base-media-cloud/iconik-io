package search

import (
	"context"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
)

// API is an interface that defines the operations that can be performed on the search endpoint.
type API interface {
	Search(ctx context.Context, path string, payload []byte) (search.ResultsDTO, error)
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

// Search searches the iconik api using the given payload.
func (s *Svc) Search(ctx context.Context, path string, payload []byte) (search.ResultsDTO, error) {
	dto, err := s.api.Search(ctx, path, payload)
	if err != nil {
		return search.ResultsDTO{}, err
	}

	if dto.Errors != nil {
		return search.ResultsDTO{}, errors.New(fmt.Sprintf("%v", dto.Errors))
	}

	return dto, nil
}

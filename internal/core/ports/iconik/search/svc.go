package search

//go:generate mockgen -source svc.go -destination=../../../../../mocks/search_mocks/svc.go -package=search_mocks

import (
	"context"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
)

// Servicer is an interface that defines the methods that a service must implement.
type Servicer interface {
	Search(ctx context.Context, path string, payload []byte) (search.ResultsDTO, error)
}

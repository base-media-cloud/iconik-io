package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
	"github.com/google/uuid"
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
		return search.ResultsDTO{}, fmt.Errorf("%v", dto.Errors)
	}

	return dto, nil
}

// ValidateAndSearchAssetID validates an asset ID, searches for it and returns it.
func (s *Svc) ValidateAndSearchAssetID(ctx context.Context, assetID, collectionID string) (search.ObjectDTO, error) {
	_, err := uuid.Parse(assetID)
	if err != nil {
		return search.ObjectDTO{}, errors.New("not a valid asset ID")
	}

	sch := search.Search{
		DocTypes:      []string{"assets", "collections"},
		Facets:        []string{"object_type", "media_type", "archive_status", "type", "format", "is_online", "approval_status"},
		IncludeFields: []string{"id", "title", "files", "in_collections", "metadata", "files.size", "media_type"},
		Sort: []search.Sort{
			{Name: "date_created", Order: "desc"},
		},
		Filter: search.Filter{
			Operator: "AND",
			Terms: []search.Term{
				{Name: "ancestor_collections", ValueIn: []string{collectionID}},
				{Name: "status", ValueIn: []string{"ACTIVE"}},
			},
		},
		FacetsFilters: []search.FacetsFilter{
			{Name: "object_type", ValueIn: []string{"assets"}},
		},
		SearchFields: []string{"title", "description", "segment_text", "file_names", "metadata", "transcription_text"},
		SearchAfter:  []interface{}{},
		Query:        assetID,
	}

	schPayload, err := json.Marshal(sch)
	if err != nil {
		return search.ObjectDTO{}, err
	}

	results, err := s.Search(ctx, iconik.SearchPath, schPayload)
	if err != nil {
		return search.ObjectDTO{}, err
	}

	if len(results.Objects) == 0 {
		return search.ObjectDTO{}, errors.New("asset not found")
	}

	return results.Objects[0], nil
}

// ValidateAndSearchFilename validates an asset filename, searches for it and returns it.
func (s *Svc) ValidateAndSearchFilename(ctx context.Context, filename, collectionID string) (search.ObjectDTO, error) {
	sch := search.Search{
		DocTypes:      []string{"assets", "collections"},
		Facets:        []string{"object_type", "media_type", "archive_status", "type", "format", "is_online", "approval_status"},
		IncludeFields: []string{"id", "title", "files", "in_collections", "metadata", "files.size", "media_type"},
		Sort: []search.Sort{
			{Name: "date_created", Order: "desc"},
		},
		Filter: search.Filter{
			Operator: "AND",
			Terms: []search.Term{
				{Name: "ancestor_collections", ValueIn: []string{collectionID}},
				{Name: "status", ValueIn: []string{"ACTIVE"}},
			},
		},
		FacetsFilters: []search.FacetsFilter{
			{Name: "object_type", ValueIn: []string{"assets"}},
		},
		SearchFields: []string{"title", "description", "segment_text", "file_names", "metadata", "transcription_text"},
		SearchAfter:  []interface{}{},
		Query:        filename,
	}

	schPayload, err := json.Marshal(sch)
	if err != nil {
		return search.ObjectDTO{}, err
	}

	results, err := s.Search(ctx, iconik.SearchPath, schPayload)
	if err != nil {
		return search.ObjectDTO{}, err
	}

	if len(results.Objects) == 0 {
		return search.ObjectDTO{}, errors.New("asset not found")
	}

	return results.Objects[0], nil
}

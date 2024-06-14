package output

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	colldomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	metadatadomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	searchdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/search"
	"strconv"
	"strings"
)

// Svc is a struct that implements the iconik servicer ports.
type Svc struct {
	collSvc     collections.Servicer
	metadataSvc metadata.Servicer
	searchSvc   search.Servicer
}

// New is a function that returns a new instance of iconik Svc struct.
func New(
	collSvc collections.Servicer,
	metadataSvc metadata.Servicer,
	searchSvc search.Servicer,
) *Svc {
	return &Svc{
		collSvc:     collSvc,
		metadataSvc: metadataSvc,
		searchSvc:   searchSvc,
	}
}

// GetMetadataView retrieves a Metadata view from the iconik API.
func (svc *Svc) GetMetadataView(ctx context.Context, viewID string) (metadatadomain.DTO, error) {
	view, err := svc.metadataSvc.GetMetadataView(ctx, iconik.MetadataPath, viewID)
	if err != nil {
		return metadatadomain.DTO{}, err
	}

	if view.Errors != nil {
		return metadatadomain.DTO{}, errors.New(fmt.Sprintf("%v", view.Errors))
	}

	return view, nil
}

// GetCollection retrieves a Collection from the iconik API.
func (svc *Svc) GetCollection(ctx context.Context, collectionID string) (colldomain.CollectionDTO, error) {
	coll, err := svc.collSvc.GetCollection(ctx, iconik.CollectionsPath, collectionID)
	if err != nil {
		return colldomain.CollectionDTO{}, err
	}

	return coll, nil
}

// ProcessPage processes each page of the iconik search results using search_after pagination.
func (svc *Svc) ProcessPage(ctx context.Context, viewFields []metadatadomain.ViewFieldDTO, collectionID string, searchAfter []interface{}, w *csv.Writer) error {
	s := searchdomain.Search{
		DocTypes:      []string{"assets", "collections"},
		Facets:        []string{"object_type", "media_type", "archive_status", "type", "format", "is_online", "approval_status"},
		IncludeFields: []string{"id", "title", "files", "in_collections", "metadata", "files.size", "media_type"},
		Sort: []searchdomain.Sort{
			{Name: "date_created", Order: "desc"},
		},
		Query: "",
		Filter: searchdomain.Filter{
			Operator: "AND",
			Terms: []searchdomain.Term{
				{Name: "ancestor_collections", ValueIn: []string{collectionID}},
				{Name: "status", ValueIn: []string{"ACTIVE"}},
			},
		},
		FacetsFilters: []searchdomain.FacetsFilter{
			{Name: "object_type", ValueIn: []string{"assets"}},
		},
		SearchFields: []string{"title", "description", "segment_text", "file_names", "metadata", "transcription_text"},
		SearchAfter:  []interface{}{},
	}

	if len(searchAfter) > 0 {
		s.SearchAfter = searchAfter
	}

	sPayload, err := json.Marshal(s)
	if err != nil {
		return err
	}

	results, err := svc.searchSvc.Search(ctx, iconik.SearchPath, sPayload)
	if err != nil {
		return err
	}

	toWrite, err := svc.FormatResultsObjects(viewFields, results.Objects)
	if err != nil {
		return err
	}

	if err = w.WriteAll(toWrite); err != nil {
		return err
	}

	if len(results.Objects) > 0 {
		lastObject := results.Objects[len(results.Objects)-1]
		searchAfterNew := lastObject.Sort
		if err = svc.ProcessPage(ctx, viewFields, collectionID, searchAfterNew, w); err != nil {
			return err
		}
	}

	return nil
}

// FormatResultsObjects formats the results of a search into a 2d slice, ready for writing.
func (svc *Svc) FormatResultsObjects(viewFields []metadatadomain.ViewFieldDTO, objs []searchdomain.ObjectDTO) ([][]string, error) {
	var metadataFile [][]string
	var csvColumnsName []string

	for _, field := range viewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
		}
	}

	numColumns := len(csvColumnsName)

	for _, object := range objs {
		row := make([]string, numColumns+4)
		row[0] = object.ID
		row[1] = "N/A"
		row[2] = "N/A"
		if len(object.Files) > 0 {
			row[1] = object.Files[0].OriginalName
			row[2] = strconv.Itoa(object.Files[0].Size)
		}
		row[3] = object.Title

		for i := 0; i < numColumns; i++ {
			metadataField := csvColumnsName[i]
			metadataValue := object.Metadata[metadataField]
			result := make([]string, len(metadataValue))

			for index, elem := range metadataValue {
				switch val := elem.(type) {
				case string:
					str := val
					if strings.HasPrefix(str, " ") {
						str = strings.TrimLeft(str, " ")
					}
					if strings.HasSuffix(str, " ") {
						str = strings.TrimRight(str, " ")
					}
					result[index] = str
				case bool:
					result[index] = fmt.Sprintf("%t", val)
				case int:
					result[index] = fmt.Sprintf("%d", val)
				case float64:
					result[index] = fmt.Sprintf("%d", int(val))
				default:
					result[index] = fmt.Sprintf("%d", val)
				}
			}

			if len(result) > 1 {
				row[i+4] = strings.Join(result, ",")
			} else {
				row[i+4] = strings.Join(result, "")
			}

		}

		metadataFile = append(metadataFile, row)
	}

	return metadataFile, nil
}

// Headers writers the headers provided by a slice of ViewFieldDTO to a 2d slice, ready for writing.
func (svc *Svc) Headers(viewFields []metadatadomain.ViewFieldDTO) [][]string {
	var metadataFile [][]string
	var csvColumnsLabel []string
	for _, field := range viewFields {
		if field.Name != "__separator__" {
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	headerRow := append([]string{"id", "original_name", "size", "title"}, csvColumnsLabel...)

	return append(metadataFile, headerRow)
}

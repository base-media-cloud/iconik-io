package output

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	collDomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	metadataDomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"strconv"
	"strings"
)

// Svc is a struct that implements the iconik servicer ports.
type Svc struct {
	collSvc     collections.Servicer
	assetSvc    assets.Servicer
	metadataSvc metadata.Servicer
}

// New is a function that returns a new instance of iconik Svc struct.
func New(
	collSvc collections.Servicer,
	assetSvc assets.Servicer,
	metadataSvc metadata.Servicer,
) *Svc {
	return &Svc{
		collSvc:     collSvc,
		assetSvc:    assetSvc,
		metadataSvc: metadataSvc,
	}
}

// GetMetadataView retrieves a Metadata view from the iconik API.
func (svc *Svc) GetMetadataView(ctx context.Context, viewID string) (metadataDomain.DTO, error) {
	view, err := svc.metadataSvc.GetMetadataView(ctx, iconik.MetadataPath, viewID)
	if err != nil {
		return metadataDomain.DTO{}, err
	}

	if view.Errors != nil {
		return metadataDomain.DTO{}, errors.New(fmt.Sprintf("%v", view.Errors))
	}

	return view, nil
}

// GetCollection retrieves a Collection from the iconik API.
func (svc *Svc) GetCollection(ctx context.Context, collectionID string) (collDomain.CollectionDTO, error) {
	coll, err := svc.collSvc.GetCollection(ctx, iconik.CollectionsPath, collectionID)
	if err != nil {
		return collDomain.CollectionDTO{}, err
	}

	return coll, nil
}

// ProcessColl takes a collection ID and recursively writes every collection
// to a csv file one collection at a time.
func (svc *Svc) ProcessColl(ctx context.Context, viewFields []metadataDomain.ViewFieldDTO, collectionID string, pageNo int, w *csv.Writer) error {
	contents, err := svc.collSvc.GetContents(ctx, iconik.CollectionsPath, collectionID, pageNo)
	if err != nil {
		return err
	}

	if contents.Errors != nil {
		return errors.New(fmt.Sprintf("%v, %v", contents.Errors, collectionID))
	}

	if err = svc.WriteCollToCSV(ctx, viewFields, contents, w); err != nil {
		return err
	}

	if contents.Pages > pageNo {
		if err = svc.ProcessColl(ctx, viewFields, collectionID, pageNo+1, w); err != nil {
			return err
		}
	}

	return nil
}

// WriteCollToCSV writes the objects from the collection to a csv file
// and will recursively call ProcessColl if another collection is found.
func (svc *Svc) WriteCollToCSV(ctx context.Context, viewFields []metadataDomain.ViewFieldDTO, contents collDomain.ContentsDTO, w *csv.Writer) error {
	var assets []collDomain.ObjectDTO

	for j := range contents.Objects {
		if contents.Objects[j].ObjectType == "collections" {
			fmt.Printf("\nfound collection %s, collection id %s", contents.Objects[j].Title, contents.Objects[j].ID)
			if err := svc.ProcessColl(ctx, viewFields, contents.Objects[j].ID, 1, w); err != nil {
				return err
			}
			continue
		}
		assets = append(assets, contents.Objects[j])
	}

	toWrite, err := svc.FormatObjects(viewFields, assets)
	if err != nil {
		return err
	}

	if err = w.WriteAll(toWrite); err != nil {
		return err
	}

	return nil
}

func (svc *Svc) FormatObjects(viewFields []metadataDomain.ViewFieldDTO, objs []collDomain.ObjectDTO) ([][]string, error) {
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

func (svc *Svc) Headers(viewFields []metadataDomain.ViewFieldDTO) [][]string {
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

package input

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	colldomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	metadatadomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/utils"
	"log"
	"os"
	"strings"
)

// Svc is a struct that implements the iconik servicer ports.
type Svc struct {
	collSvc     collections.Servicer
	assetSvc    assets.Servicer
	metadataSvc metadata.Servicer
}

// New is a function that returns a new instance of the iconik Svc struct.
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

// GetCollectionObjects gets all the results from a collection and returns the full object list.
func (svc *Svc) GetCollectionObjects(ctx context.Context, collectionID string, pageNo int, objects []colldomain.ObjectDTO) ([]colldomain.ObjectDTO, error) {
	coll, err := svc.collSvc.GetContents(ctx, iconik.CollectionsPath, collectionID, pageNo)
	if err != nil {
		return nil, err
	}

	if coll.Errors != nil {
		return nil, errors.New(fmt.Sprintf("%v", coll.Errors))
	}

	objects = append(objects, coll.Objects...)

	if coll.Pages > 1 && coll.Pages > pageNo {
		objs, err := svc.GetCollectionObjects(ctx, collectionID, pageNo+1, objects)
		if err != nil {
			return nil, err
		}

		return objs, nil
	}

	return objects, nil
}

// ProcessObjects takes a slice of objects and returns assets only.
func (svc *Svc) ProcessObjects(ctx context.Context, assets, objects []colldomain.ObjectDTO, assetsMap, collectionsMap map[string]struct{}) ([]colldomain.ObjectDTO, error) {
	for _, o := range objects {
		if o.ObjectType == "assets" {
			if _, exists := assetsMap[o.ID]; !exists {
				assets = append(assets, o)
				assetsMap[o.ID] = struct{}{}
			}
		} else if o.ObjectType == "collections" {
			if _, exists := collectionsMap[o.ID]; !exists {
				fmt.Println()
				fmt.Printf("found collection %s, traversing:\n", o.Title)
				var err error
				var objs []colldomain.ObjectDTO
				objs, err = svc.GetCollectionObjects(ctx, o.ID, 1, objs)
				if err != nil {
					fmt.Println("Error fetching data for collection with ID", o.ID, "Error:", err)
					continue
				}

				collectionsMap[o.ID] = struct{}{}
				a, err := svc.ProcessObjects(ctx, assets, objs, assetsMap, collectionsMap)
				if err != nil {
					return nil, err
				}
				return a, nil
			}
		}
	}
	return assets, nil
}

// ReadCSVFile reads a CSV file and returns it as a 2D slice.
func (svc *Svc) ReadCSVFile(appCfg *config.App) ([][]string, error) {
	csvFile, err := os.Open(appCfg.Input)
	if err != nil {
		return nil, errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData, nil
}

// UpdateIconik reads a 2D slice, verifies it, and uploads the data to the Iconik API.
func (svc *Svc) UpdateIconik(ctx context.Context, viewFields []metadatadomain.ViewFieldDTO, assets []colldomain.ObjectDTO, metadataFile [][]string, cfg *config.App) error {
	csvHeaders := metadataFile[0]
	if csvHeaders[0] != "id" || csvHeaders[1] != "original_name" || csvHeaders[2] != "size" || csvHeaders[3] != "title" {
		fmt.Println(csvHeaders)
		return errors.New("CSV file not properly formatted for Iconik")
	}

	matchingData, nonMatchingHeaders, err := utils.MatchCSVtoAPI(viewFields, metadataFile)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		fmt.Println("Some columns from the file provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
		fmt.Println()
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	matchingFileHeaderNames := matchingData[0]
	matchingFileHeaderLabels := matchingData[1]

	csvFilesToUpdate := len(matchingData) - 2
	fmt.Println("Amount of files to update:", csvFilesToUpdate)

	notAdded := make(map[string]string)

	for index := 2; index < len(matchingData); index++ {
		row := matchingData[index]

		assetID := row[0]
		origName := row[1]
		title := row[3]

		errAssetID := svc.assetSvc.ValidateAsset(ctx, assetID)
		for _, asset := range assets {
			for _, file := range asset.Files {
				if asset.ID == assetID {
					origName = file.OriginalName
				}
			}
		}

		aID, errFilename := utils.ValidateFilename(assets, origName)
		if errFilename == nil {
			assetID = aID
		}

		if errAssetID != nil && errFilename != nil {
			log.Printf("%s & %s, skipping\n", errAssetID, errFilename)
			notAdded[assetID] = origName
			continue
		}

		metadataValues := metadatadomain.Values{
			MetadataValues: map[string]struct {
				FieldValues []metadatadomain.FieldValue `json:"field_values"`
			}(make(map[string]struct {
				FieldValues []metadatadomain.FieldValue
			})),
		}

		for count := 4; count < len(row); count++ {
			headerName := matchingFileHeaderNames[count]
			headerLabel := matchingFileHeaderLabels[count]
			fieldValueSlice := make([]metadatadomain.FieldValue, 0)

			valueArr := strings.Split(row[count], ",")
			if len(valueArr) > 1 {
				for _, val := range valueArr {

					err = utils.ValidateSchema(headerLabel, val)
					if err != nil {
						return err
					}

					fieldValueSlice = append(fieldValueSlice, metadatadomain.FieldValue{Value: val})
				}
				metadataValues.MetadataValues[headerName] = struct {
					FieldValues []metadatadomain.FieldValue `json:"field_values"`
				}{
					FieldValues: fieldValueSlice,
				}
			}
		}

		assetPayload, err := json.Marshal(map[string]string{"title": title})
		if err != nil {
			return errors.New("error marshaling JSON")
		}

		_, err = svc.assetSvc.UpdateAsset(ctx, iconik.AssetsPath, assetID, assetPayload)
		if err != nil {
			log.Println("Error updating title for asset ", assetID)
			return err
		}

		metadataPayload, err := json.Marshal(metadataValues)
		if err != nil {
			return errors.New("error marshaling JSON")
		}

		_, err = svc.metadataSvc.UpdateMetadataInAsset(ctx, iconik.AssetsPath, cfg.ViewID, assetID, metadataPayload)
		if err != nil {
			log.Println("Error updating metadata for asset ", assetID)
			return err
		}
	}

	fmt.Printf("Assets successfully updated: %d of %d\n", csvFilesToUpdate-len(notAdded), csvFilesToUpdate)
	if len(notAdded) > 0 {
		fmt.Println("Some assets failed to update:")
		for assetID, origName := range notAdded {
			fmt.Printf("Asset ID: %s, Original filename: %s", assetID, origName)
		}
	}

	return nil
}

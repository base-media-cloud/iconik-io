package input

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	metadatadomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/search"
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
	searchSvc   search.Servicer
}

// New is a function that returns a new instance of the iconik Svc struct.
func New(
	collSvc collections.Servicer,
	assetSvc assets.Servicer,
	metadataSvc metadata.Servicer,
	searchSvc search.Servicer,
) *Svc {
	return &Svc{
		collSvc:     collSvc,
		assetSvc:    assetSvc,
		metadataSvc: metadataSvc,
		searchSvc:   searchSvc,
	}
}

func (svc *Svc) ProcessAssets(ctx context.Context, csvData [][]string, collectionID, viewID string) (map[string]bool, error) {
	matchingFileHeaderNames := csvData[0]
	matchingFileHeaderLabels := csvData[1]
	notAdded := make(map[string]bool)

	for i := 2; i < len(csvData); i++ {
		row := csvData[i]
		assetID := row[0]
		origName := row[1]
		title := row[3]

		_, errAssetID := svc.searchSvc.ValidateAndSearchAssetID(ctx, assetID, collectionID)
		if errAssetID != nil {
			result, errFilename := svc.searchSvc.ValidateAndSearchFilename(ctx, origName, collectionID)
			if errFilename != nil {
				log.Printf("%s & %s for %s, skipping\n", errAssetID, errFilename, title)
				notAdded[assetID] = true
				continue
			}
			assetID = result.ID
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
			if len(valueArr) > 0 {
				for _, val := range valueArr {
					err := utils.ValidateSchema(headerLabel, val)
					if err != nil {
						return nil, err
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
			return nil, errors.New("error marshaling JSON")
		}

		_, err = svc.assetSvc.UpdateAsset(ctx, iconik.AssetsPath, assetID, assetPayload)
		if err != nil {
			log.Println("Error updating title for asset ", assetID)
			return nil, err
		}

		metadataPayload, err := json.Marshal(metadataValues)
		if err != nil {
			return nil, errors.New("error marshaling JSON")
		}

		_, err = svc.metadataSvc.UpdateMetadataInAsset(ctx, iconik.MetadataAssetsPath, viewID, assetID, metadataPayload)
		if err != nil {
			return nil, fmt.Errorf("error updating metadata for asset %v", assetID)
		}

	}

	return notAdded, nil
}

// GetMetadataView retrieves a Metadata view from the iconik API.
func (svc *Svc) GetMetadataView(ctx context.Context, viewID string) (metadatadomain.DTO, error) {
	view, err := svc.metadataSvc.GetMetadataView(ctx, iconik.MetadataViewPath, viewID)
	if err != nil {
		return metadatadomain.DTO{}, err
	}

	if view.Errors != nil {
		return metadatadomain.DTO{}, fmt.Errorf("%v", view.Errors)
	}

	return view, nil
}

// ReadCSVFile reads a CSV file and returns it as a 2D slice.
func (svc *Svc) ReadCSVFile(appCfg *config.App) ([][]string, error) {
	csvFile, err := os.Open(appCfg.Input)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData, nil
}

// MatchCSVtoView takes a csv as a 2d slice, and checks its fields against the inputted view field from iconik.
func (svc *Svc) MatchCSVtoView(viewFields []metadatadomain.ViewFieldDTO, csvData [][]string) ([][]string, []string, error) {
	csvHeaderLabels := csvData[0]

	var matchingIconikHeaderNames []string
	var matchingIconikHeaderLabels []string
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "id")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "original_name")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "size")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "title")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "id")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "original_name")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "size")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "title")

	var nonMatchingHeaders []string

	for index, csvHeaderLabel := range csvHeaderLabels {
		if index > 3 {
			found := false
			for _, viewField := range viewFields {
				if csvHeaderLabel == viewField.Label {
					matchingIconikHeaderNames = append(matchingIconikHeaderNames, viewField.Name)
					matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, viewField.Label)
					found = true
					break
				}
			}
			if !found {
				nonMatchingHeaders = append(nonMatchingHeaders, csvHeaderLabel)
			}
		}
	}

	var matchingValues [][]string
	matchingValues = append(matchingValues, matchingIconikHeaderNames)
	matchingValues = append(matchingValues, matchingIconikHeaderLabels)

	for j := 1; j < len(csvData); j++ {
		row := csvData[j]
		var matchingRow []string
		for k, csvHeaderLabel := range csvHeaderLabels {
			if contains(matchingIconikHeaderLabels, csvHeaderLabel) {
				matchingRow = append(matchingRow, row[k])
			}
		}
		matchingValues = append(matchingValues, matchingRow)
	}

	return matchingValues, nonMatchingHeaders, nil

}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

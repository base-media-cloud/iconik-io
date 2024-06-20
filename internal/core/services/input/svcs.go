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
	searchdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/search"
	"github.com/base-media-cloud/pd-iconik-io-rd/utils"
	"github.com/rs/zerolog"
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

func (svc *Svc) ProcessAssets(ctx context.Context, csvData [][]string, collectionID, viewID string) (map[string]string, error) {
	matchingFileHeaderNames := csvData[0]
	matchingFileHeaderLabels := csvData[1]
	notAdded := make(map[string]string)

	for i := 2; i < len(csvData); i++ {
		row := csvData[i]
		assetID := row[0]
		origName := row[1]
		title := row[3]

		_, errAssetID := svc.searchSvc.ValidateAndSearchAssetID(ctx, assetID, collectionID)
		if errAssetID != nil {
			_, errFilename := svc.searchSvc.ValidateAndSearchFilename(ctx, origName, collectionID)
			if errAssetID != nil && errFilename != nil {
				log.Printf("%s & %s, skipping\n", errAssetID, errFilename)
				notAdded[assetID] = origName
				continue
			}
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
			if !isBlankStringArray(valueArr) {
				for _, val := range valueArr {

					err = SchemaValidator(headerLabel, val)
					if err != nil {
						return err
					}

					fieldValueSlice = append(fieldValueSlice, FieldValue{Value: val})
				}
				csvMetadata.MetadataValuesStruct.MetadataValues[headerName] = struct {
					FieldValues []FieldValue `json:"field_values"`
				}{
					FieldValues: fieldValueSlice,
				}
			} else {
				continue
			}
			// if len(valueArr) > 1 {
			// 	for _, val := range valueArr {
			// 		err := utils.ValidateSchema(headerLabel, val)
			// 		if err != nil {
			// 			return nil, err
			// 		}
			//
			// 		fieldValueSlice = append(fieldValueSlice, metadatadomain.FieldValue{Value: val})
			// 	}
			// 	metadataValues.MetadataValues[headerName] = struct {
			// 		FieldValues []metadatadomain.FieldValue `json:"field_values"`
			// 	}{
			// 		FieldValues: fieldValueSlice,
			// 	}
			// }
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

		fmt.Println(string(metadataPayload))

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
		return metadatadomain.DTO{}, errors.New(fmt.Sprintf("%v", view.Errors))
	}

	return view, nil
}

// // GetCollectionObjects gets all the results from a collection and returns the full object list.
// func (svc *Svc) GetCollectionObjects(ctx context.Context, collectionID string, pageNo int, objects []colldomain.ObjectDTO) ([]colldomain.ObjectDTO, error) {
// 	coll, err := svc.collSvc.GetContents(ctx, iconik.CollectionsPath, collectionID, pageNo)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if coll.Errors != nil {
// 		return nil, errors.New(fmt.Sprintf("%v", coll.Errors))
// 	}
//
// 	objects = append(objects, coll.Objects...)
//
// 	if coll.Pages > 1 && coll.Pages > pageNo {
// 		objs, err := svc.GetCollectionObjects(ctx, collectionID, pageNo+1, objects)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		return objs, nil
// 	}
//
// 	return objects, nil
// }

// ProcessPage processes each page of the iconik search results using search_after pagination.
func (svc *Svc) ProcessPage(ctx context.Context, csvFilesToUpdate int, matchingData [][]string, viewFields []metadatadomain.ViewFieldDTO, collectionID, viewID string, searchAfter []interface{}) error {
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

	err = svc.UpdateIconik(ctx, csvFilesToUpdate, viewFields, viewID, results.Objects, matchingData)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to update iconik")
		return err
	}

	if len(results.Objects) > 0 {
		lastObject := results.Objects[len(results.Objects)-1]
		searchAfterNew := lastObject.Sort
		if err = svc.ProcessPage(ctx, csvFilesToUpdate, matchingData, viewFields, collectionID, viewID, searchAfterNew); err != nil {
			return err
		}
	}

	return nil
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

// UpdateIconik reads a 2D slice, verifies it, and uploads the data to the Iconik API.
func (svc *Svc) UpdateIconik(ctx context.Context, csvFilesToUpdate int, viewFields []metadatadomain.ViewFieldDTO, viewID string, assets []searchdomain.ObjectDTO, matchingData [][]string) error {
	// csvHeaders := metadataFile[0]
	// if csvHeaders[0] != "id" || csvHeaders[1] != "original_name" || csvHeaders[2] != "size" || csvHeaders[3] != "title" {
	// 	fmt.Println(csvHeaders)
	// 	return errors.New("CSV file not properly formatted for Iconik")
	// }
	//
	// matchingData, nonMatchingHeaders, err := utils.MatchCSVtoAPI(viewFields, metadataFile)
	// if err != nil {
	// 	return err
	// }
	//
	// if len(nonMatchingHeaders) > 0 {
	// 	fmt.Println("Some columns from the file provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
	// 	fmt.Println()
	// 	for _, nonMatchingHeader := range nonMatchingHeaders {
	// 		fmt.Println(nonMatchingHeader)
	// 	}
	// }
	//
	matchingFileHeaderNames := matchingData[0]
	matchingFileHeaderLabels := matchingData[1]
	//
	// csvFilesToUpdate := len(matchingData) - 2
	// fmt.Println("Amount of files to update:", csvFilesToUpdate)

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
					err := utils.ValidateSchema(headerLabel, val)
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

		_, err = svc.metadataSvc.UpdateMetadataInAsset(ctx, iconik.AssetsPath, viewID, assetID, metadataPayload)
		if err != nil {
			log.Println("Error updating metadata for asset ", assetID)
			return err
		}
	}

	return nil
}

// // ProcessObjects takes a slice of objects and returns assets only.
// func (svc *Svc) ProcessObjects(ctx context.Context, assets, objects []colldomain.ObjectDTO, assetsMap, collectionsMap map[string]struct{}) ([]colldomain.ObjectDTO, error) {
// 	for _, o := range objects {
// 		if o.ObjectType == "assets" {
// 			if _, exists := assetsMap[o.ID]; !exists {
// 				assets = append(assets, o)
// 				assetsMap[o.ID] = struct{}{}
// 			}
// 		} else if o.ObjectType == "collections" {
// 			if _, exists := collectionsMap[o.ID]; !exists {
// 				fmt.Println()
// 				fmt.Printf("found collection %s, traversing:\n", o.Title)
// 				var err error
// 				var objs []colldomain.ObjectDTO
// 				objs, err = svc.GetCollectionObjects(ctx, o.ID, 1, objs)
// 				if err != nil {
// 					fmt.Println("Error fetching data for collection with ID", o.ID, "Error:", err)
// 					continue
// 				}
//
// 				collectionsMap[o.ID] = struct{}{}
// 				a, err := svc.ProcessObjects(ctx, assets, objs, assetsMap, collectionsMap)
// 				if err != nil {
// 					return nil, err
// 				}
// 				return a, nil
// 			}
// 		}
// 	}
// 	return assets, nil
// }

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

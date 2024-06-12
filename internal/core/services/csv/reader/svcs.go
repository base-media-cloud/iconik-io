package reader

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/ports/iconik/validate"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	csvdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/csv"
	"github.com/base-media-cloud/pd-iconik-io-rd/utils"
)

// Svc is a struct that implements the iconik servicer ports.
type Svc struct {
	collSvc  collections.Servicer
	assetSvc assets.Servicer
	val      validate.Validator
}

// New is a function that returns a new instance of iconik Svc struct.
func New(
	collSvc collections.Servicer,
	assetSvc assets.Servicer,
	val validate.Validator,
) *Svc {
	return &Svc{
		collSvc:  collSvc,
		assetSvc: assetSvc,
		val:      val,
	}
}

func (svc *Svc) ReadCSVFile(iconikCfg *config.Iconik) ([][]string, error) {
	csvFile, err := os.Open(iconikCfg.Input)
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
func (svc *Svc) UpdateIconik(viewFields []*metadata.ViewField, metadataFile [][]string, iconikCfg *config.Iconik) error {
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

	var c csvdomain.CSV

	c.CSVFilesToUpdate = len(matchingData) - 2
	fmt.Println("Amount of files to update:", c.CSVFilesToUpdate)

	for index := 2; index < len(matchingData); index++ {
		row := matchingData[index]

		csvMetadata := csvdomain.CSVMetadata{
			Added: false,
			IDStruct: csvdomain.IDStruct{
				ID: row[0],
			},
			OriginalNameStruct: csvdomain.OriginalNameStruct{
				OriginalName: row[1],
			},
			SizeStruct: csvdomain.SizeStruct{
				Size: row[2],
			},
			TitleStruct: csvdomain.TitleStruct{
				Title: row[3],
			},
			MetadataValuesStruct: csvdomain.MetadataValuesStruct{
				MetadataValues: make(map[string]struct {
					FieldValues []csvdomain.FieldValue `json:"field_values"`
				}),
			},
		}

		c.CSVMetadata = append(c.CSVMetadata, &csvMetadata)

		errAssetID := svc.val.AssetID(objects, csvMetadata, ctx)
		errFilename := svc.val.Filename(objects, csvMetadata)

		if errAssetID != nil && errFilename != nil {
			log.Printf("%s & %s, skipping\n", errAssetID, errFilename)
			continue
		}
		csvMetadata.Added = true

		for count := 4; count < len(row); count++ {
			headerName := matchingFileHeaderNames[count]
			headerLabel := matchingFileHeaderLabels[count]
			fieldValueSlice := make([]csvdomain.FieldValue, 0)

			valueArr := strings.Split(row[count], ",")
			if !utils.IsBlankStringArray(valueArr) {
				for _, val := range valueArr {

					err = svc.val.Schema(headerLabel, val)
					if err != nil {
						return err
					}

					fieldValueSlice = append(fieldValueSlice, csvdomain.FieldValue{Value: val})
				}
				csvMetadata.MetadataValuesStruct.MetadataValues[headerName] = struct {
					FieldValues []csvdomain.FieldValue `json:"field_values"`
				}{
					FieldValues: fieldValueSlice,
				}
			} else {
				continue
			}
		}

		assetPayload, err := json.Marshal(csvMetadata.TitleStruct)
		if err != nil {
			return errors.New("error marshaling JSON")
		}

		_, err = svc.assetSvc.UpdateAsset(ctx, iconik.AssetsPath, csvMetadata.IDStruct.ID, assetPayload)
		if err != nil {
			log.Println("Error updating title name for asset ", csvMetadata.IDStruct.ID)
			return err
		}

		metadataPayload, err := json.Marshal(csvMetadata.MetadataValuesStruct)
		if err != nil {
			return errors.New("error marshaling JSON")
		}

		_, err = svc.api.UpdateMetadataInAsset(ctx, iconik.AssetsPath, iconikCfg.ViewID, csvMetadata.IDStruct.ID, metadataPayload)
		if err != nil {
			log.Println("Error updating metadata for asset ", csvMetadata.IDStruct.ID)
			return err
		}
	}

	fmt.Println()
	log.Println("Assets successfully updated:")
	var countSuccess int
	for _, csvMetadata := range c.CSVMetadata {
		if csvMetadata.Added {
			countSuccess++
		}
	}
	fmt.Printf("%d of %d", countSuccess, c.CSVFilesToUpdate)

	fmt.Println()
	log.Println("Assets that failed to update:")
	var countFailed int
	for _, csvMetadata := range c.CSVMetadata {
		if !csvMetadata.Added {
			countFailed++
			log.Printf("%s (Title: %s, Original filename: %s)", csvMetadata.IDStruct.ID, csvMetadata.TitleStruct.Title, csvMetadata.OriginalNameStruct.OriginalName)
		}
	}
	fmt.Printf("%d of %d\n", countFailed, c.CSVFilesToUpdate)

	return nil
}

// GetCollection gets all the results from a collection and return the full object list with metadata.
func (i *Iconik) GetCollection(collectionID string, pageNo int) error {
	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Path, collectionID, "/contents/")
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme
	queryParams := u.Query()
	queryParams.Set("per_page", "500")
	queryParams.Set("page", strconv.Itoa(pageNo))
	u.RawQuery = queryParams.Encode()

	_, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}

	switch {
	default:
		err = json.Unmarshal(resBody, &i.IconikClient.Collection)
		if err != nil {
			return err
		}
	case i.IconikClient.Collection != nil:
		var c *Collection
		err = json.Unmarshal(resBody, &c)
		if err != nil {
			return err
		}
		i.IconikClient.Collection.Objects = append(i.IconikClient.Collection.Objects, c.Objects...)
	}

	if i.IconikClient.Collection.Errors != nil {
		return errors.New(fmt.Sprintf("%v", i.IconikClient.Collection.Errors))
	}

	if i.IconikClient.Collection.Pages > 1 && i.IconikClient.Collection.Pages > pageNo {
		if err := i.GetCollection(collectionID, pageNo+1); err != nil {
			return err
		}
	}

	return nil
}

func (i *Iconik) ProcessObjects(c *Collection, assetsMap, collectionsMap map[string]struct{}) error {
	for _, o := range c.Objects {
		if o.ObjectType == "assets" {
			if _, exists := assetsMap[o.ID]; !exists {
				i.IconikClient.Assets = append(i.IconikClient.Assets, o)
				assetsMap[o.ID] = struct{}{}
			}
		} else if o.ObjectType == "collections" {
			if _, exists := collectionsMap[o.ID]; !exists {
				fmt.Println()
				fmt.Printf("found collection %s, traversing:\n", o.Title)
				err := i.GetCollection(o.ID, 1)
				if err != nil {
					fmt.Println("Error fetching data for collection with ID", o.ID, "Error:", err)
					continue
				}
				collectionsMap[o.ID] = struct{}{}
				if err := i.ProcessObjects(i.IconikClient.Collection, assetsMap, collectionsMap); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

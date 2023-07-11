package iconikio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

// ReadCSVFile reads and validates the CSV file provided.
func (i *Iconik) ReadCSVFile() error {

	csvFile, err := os.Open(i.IconikClient.Config.Input)
	if err != nil {
		return errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	csvData, err := csvReader.ReadAll()
	if err != nil {
		return errors.New("error reading CSV file")
	}

	csvHeaders := csvData[0]

	if csvHeaders[0] != "id" || csvHeaders[1] != "title" {
		return errors.New("CSV file not properly formatted for Iconik")
	}

	_, _, err = i.GetCSVColumnsFromView()
	if err != nil {
		return err
	}
	
	matchingCSV, nonMatchingHeaders, err := i.matchCSVtoAPI(csvData)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		fmt.Println("Some columns from the CSV provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
		fmt.Println()
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	matchingCSVHeaderNames := matchingCSV[0]
	matchingCSVHeaderLabels := matchingCSV[1]

	for index, row := range matchingCSV {
		if index > 1 {

			title := make(map[string]string)
			metadata := make(map[string]interface{})
			metadataValues := make(map[string]interface{})

			for count, value := range row {
				if count == 0 {
					// it's the asset id
					// check asset id is valid
					_, err := uuid.Parse(value)
					if err != nil {
						return errors.New("not a valid asset ID")
					}
	
					_, err = i.CheckAssetbyID(value)
					if err != nil {
						return fmt.Errorf("error %s", err)
					}
	
					code, err := i.CheckAssetExistInCollection(value)
					if err != nil {
						return err
					}
					if code == http.StatusOK {
						continue
					} else {
						return errors.New("asset does not exist in given Collection ID")
					}
	
				} else if count == 1 {
					// it's the title of the asset
					title["title"] = value
				} else if count > 1 {
					// this is where the rest of the headers start
					headerName := matchingCSVHeaderNames[count]
					headerLabel := matchingCSVHeaderLabels[count]
	
					if _, ok := metadataValues[headerName]; !ok {
						metadataValues[headerName] = map[string]interface{}{
							"field_values": []map[string]interface{}{},
						}
					}
	
					valueArr := strings.Split(value, ",")
					if len(valueArr) > 0 {
						for _, val := range valueArr {
	
							_, val, err = SchemaValidator(headerLabel, val)
							if err != nil {
								return err
							}
	
							fieldValue := map[string]interface{}{
								"value": val,
							}
	
							if val != "" {
								fieldValues := metadataValues[headerName].(map[string]interface{})["field_values"].([]map[string]interface{})
								fieldValues = append(fieldValues, fieldValue)
								metadataValues[headerName].(map[string]interface{})["field_values"] = fieldValues
							} else {
								delete(metadataValues, headerName)
							}
						}
					}
				}
			}

			err = i.updateTitle(row[0], title)
			if err != nil {
				return err
			}
	
			metadata["metadata_values"] = metadataValues
	
			err = i.updateMetadata(row[0], metadata)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// updateTitle updates the title for the given asset ID.
func (i *Iconik) updateTitle(assetID string, title map[string]string) error {

	requestBody, err := json.Marshal(title)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	// uri := i.IconikClient.Config.IconikURL + "/API/assets/v1/assets/" + assetID

	uri, err := i.joinURL("asset", assetID, 1)
	if err != nil {
		return err
	}

	res, _, err := i.getResponseBody("PATCH", uri.String(), bytes.NewBuffer(requestBody))

	if res.StatusCode == 200 {
		log.Println("Successfully updated title name for asset ", assetID)
	} else {
		log.Println("Error updating title name for asset ", assetID)
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

// updateMetadata updates the metadata for the given asset ID.
func (i *Iconik) updateMetadata(assetID string, metadata map[string]interface{}) error {

	requestBody, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	// uri := i.IconikClient.Config.IconikURL + "/API/metadata/v1/assets/" + assetID + "/views/" + i.IconikClient.Config.ViewID + "/"

	uri, err := i.joinURL("metadataView", assetID, 1)
	if err != nil {
		return err
	}

	res, _, err := i.getResponseBody("PUT", uri.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Println("Successfully updated metadata for asset ", assetID)
	} else {
		log.Println("Error updating metadata for asset ", assetID)
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

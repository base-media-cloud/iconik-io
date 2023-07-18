package iconikio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
)

// ReadCSVFile reads and validates the CSV file provided.
func (i *Iconik) ReadCSVFile() error {

	// open the provided CSV file
	csvFile, err := os.Open(i.IconikClient.Config.Input)
	if err != nil {
		return errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	// create a new reader
	csvReader := csv.NewReader(csvFile)

	// read all the CSV data into a 2D slice we can work with
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return errors.New("error reading CSV file")
	}

	// the first row of the 2D slice will be the header row
	csvHeaders := csvData[0]

	// we then validate that the schema for the header row is correct
	if csvHeaders[0] != "id" || csvHeaders[1] != "title" {
		return errors.New("CSV file not properly formatted for Iconik")
	}

	// get the slimmed down 2D slice that contains only the columns matched to the given metadata view
	matchingCSV, nonMatchingHeaders, err := i.matchCSVtoAPI(csvData)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		// log to the user if there were any columns in the provided CSV that have been left out
		fmt.Println("Some columns from the CSV provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
		fmt.Println()
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	// get the CSV Header Names row
	matchingCSVHeaderNames := matchingCSV[0]
	// get the CSV Header Labels row
	matchingCSVHeaderLabels := matchingCSV[1]

	// range over the remaining rows, which will be the assets
	for index, row := range matchingCSV {
		if index > 1 {

			// make maps to store the values. these will then be used to update the Iconik API
			title := make(map[string]string)
			metadata := make(map[string]interface{})
			metadataValues := make(map[string]interface{})

			// range over the values in the row
			for count, value := range row {
				if count == 0 {
					// it's the asset id
					// check asset id is valid
					_, err := uuid.Parse(value)
					if err != nil {
						return errors.New("not a valid asset ID")
					}

					// check asset id exists on Iconik
					_, err = i.CheckAssetbyID(value)
					if err != nil {
						return fmt.Errorf("error %s", err)
					}

					// check asset id exists in given collection id
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

					// if there are words separated by spaces inside the field, split with a comma
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

			// update the title name for this particular asset/row
			err = i.updateTitle(row[0], title)
			if err != nil {
				return err
			}

			metadata["metadata_values"] = metadataValues

			// update the remaining metadata for this particular asset/row
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

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Asset.Patch.Path, assetID)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, _, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Asset.Patch.Method, u.String(), bytes.NewBuffer(requestBody))

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

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Path, assetID, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Path2)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Method, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Println("Successfully updated metadata for asset ", assetID)
	} else {
		log.Println("Error updating metadata for asset ", assetID)
		log.Println(string(resBody))
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

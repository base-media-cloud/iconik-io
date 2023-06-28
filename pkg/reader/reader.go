package reader

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/validate"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/web"
)

func ReadCSVFile(cfg *config.Conf, log *zap.SugaredLogger) error {

	// Open CSV file
	csvFile, err := os.Open(cfg.Input)
	if err != nil {
		return errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	// Create reader
	csvReader := csv.NewReader(csvFile)

	// Read the first row of the file to get the field names
	fields, err := csvReader.Read()
	if err != nil {
		return errors.New("error reading CSV file")
	}

	// Check for properly formatted headers
	if fields[0] != "id" || fields[1] != "title" {
		return errors.New("CSV file not properly formatted for Iconik")
	}

	// Loop through each row of the CSV
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New("error reading CSV file")
		}

		// Create our maps for titles and metadata
		title := make(map[string]string)
		metadata := make(map[string]interface{})
		metadataValues := make(map[string]interface{})

		// Loop over each value in the row
		for count, value := range row {
			// First column is always our asset id
			if count == 0 {
				// UUID validation
				_, err := uuid.Parse(value)
				if err != nil {
					return errors.New("not a valid asset ID")
				}

				// Check asset exists on Iconik servers
				_, err = validate.CheckAssetbyID(value, cfg, log)
				if err != nil {
					return fmt.Errorf("error %s", err)
				}

				// Check asset is in collection provided
				code, err := validate.CheckAssetExistInCollection(value, cfg, log)
				if err != nil {
					return err
				}
				if code == http.StatusOK {
					continue
				} else {
					return errors.New("asset does not exist in given Collection ID")
				}

			} else if count == 1 {
				// Second column is always our title
				title["title"] = value
			} else if count > 1 {

				// Columns after that are our metadata and variable in length
				header := fields[count]

				// Check if the header exists in metadataValues
				if _, ok := metadataValues[header]; !ok {
					// Create a new field values slice
					metadataValues[header] = map[string]interface{}{
						"field_values": []map[string]interface{}{},
					}
				}

				// Turn all the field values into an array, even if there is only one
				valueArr := strings.Split(value, ",")

				if len(valueArr) > 0 {
					// Range over the array of substrings
					for _, val := range valueArr {

						// Validate the schema
						header, val, err = validate.SchemaValidator(header, val)
						if err != nil {
							return err
						}

						// Create a field value map
						fieldValue := map[string]interface{}{
							"value": val,
						}

						// Check if there is even anything in the column
						if val != "" {
							// Append the field value to the slice
							fieldValues := metadataValues[header].(map[string]interface{})["field_values"].([]map[string]interface{})
							fieldValues = append(fieldValues, fieldValue)
							metadataValues[header].(map[string]interface{})["field_values"] = fieldValues
						} else {
							// If metadata not in column, remove the key
							delete(metadataValues, header)
						}
					}
				}
			}
		}
		// Update the title
		err = updateTitle(cfg, row[0], title, log)
		if err != nil {
			return err
		}

		// Assign the metadata values to the metadata map
		metadata["metadata_values"] = metadataValues

		// Update the metadata
		err = updateMetadata(cfg, row[0], metadata, log)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateTitle updates the title for the given asset ID.
func updateTitle(cfg *config.Conf, assetID string, title map[string]string, log *zap.SugaredLogger) error {

	requestBody, err := json.Marshal(title)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	uri := cfg.IconikURL + "/API/assets/v1/assets/" + assetID
	res, _, err := web.GetResponseBody("PATCH", uri, bytes.NewBuffer(requestBody), cfg, log)

	if res.StatusCode == 200 {
		log.Info("Successfully updated title name for asset ", assetID)
	} else {
		log.Info("Error updating title name for asset ", assetID)
		log.Infow(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

// updateMetadata updates the metadata for the given asset ID.
func updateMetadata(cfg *config.Conf, assetID string, metadata map[string]interface{}, log *zap.SugaredLogger) error {

	requestBody, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	uri := cfg.IconikURL + "/API/metadata/v1/assets/" + assetID + "/views/" + cfg.ViewID + "/"
	res, _, err := web.GetResponseBody("PUT", uri, bytes.NewBuffer(requestBody), cfg, log)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Info("Successfully updated metadata for asset ", assetID)
	} else {
		log.Info("Error updating metadata for asset ", assetID)
		log.Infow(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

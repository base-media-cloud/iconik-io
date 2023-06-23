package reader

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/assets"
)

func ReadCSVFile(cfg *config.Conf) error {

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
				_, err = assets.GetAssetbyID(value, cfg)
				if err != nil {
					return fmt.Errorf("Error %s", err)
				}

				// Check asset is in collection provided
				bool, err := assets.DoesAssetExistInCollection(value, cfg)
				if err != nil {
					return err
				}
				if bool {
					continue
				} else {
					return errors.New("Asset does not exist in given Collection ID")
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
						header, val, err = schemaValidator(header, val)
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
		err = updateTitle(cfg, row[0], title)
		if err != nil {
			return err
		}

		// Assign the metadata values to the metadata map
		metadata["metadata_values"] = metadataValues

		// Update the metadata
		err = updateMetadata(cfg, row[0], metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateTitle updates the title for the given asset ID.
func updateTitle(cfg *config.Conf, assetID string, title map[string]string) error {
	uri := cfg.IconikURL + "/API/assets/v1/assets/" + assetID
	log.Println(uri)
	method := "PATCH"

	requestBody, err := json.Marshal(title)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Println("Successfully updated title name for asset", assetID)
	} else {
		log.Println("Error updating title name for asset", assetID)
		log.Println(res.StatusCode)
		return err
	}

	return nil
}

// updateMetadata updates the metadata for the given asset ID.
func updateMetadata(cfg *config.Conf, assetID string, metadata map[string]interface{}) error {

	uri := cfg.IconikURL + "/API/metadata/v1/assets/" + assetID + "/views/" + cfg.ViewID + "/"
	log.Println(uri)
	method := "PUT"

	requestBody, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Println("Successfully updated metadata for asset", assetID)
	} else {
		log.Println("Error updating metadata for asset", assetID)
		log.Println(res.StatusCode)
		return err
	}

	return nil
}

func schemaValidator(header, val string) (string, string, error) {
	// Schema validation for boolean fields
	if header == "Signedoff" || header == "win_Archived" || header == "ShareNo" || header == "bmc_sapProductAssetOnly" {
		if val == "TRUE" {
			val = "true"
		} else if val == "FALSE" {
			val = "false"
		} else if val == "true" || val == "false" {
		} else {
			return header, val, fmt.Errorf("For %s the value must either be set to true or false. The value is currently set to: %s", header, val)
		}
	}

	if header == "pdTest_FrameRate" || header == "pdTest_AudioFrameRate" {
		if val == "23.976" || val == "23.98" || val == "24" || val == "25" || val == "29.97" || val == "30" || val == "50" || val == "59.94" || val == "60" {
		} else {
			return header, val, fmt.Errorf("For %s the value must either be set to 23.976, 23.98, 24, 25, 29.97, 30, 50, 59.94 or 60. The value is currently set to: %s", header, val)
		}
	}

	if header == "pdTest_FrameRateMode" {
		if val == "Constant" || val == "Variable" {
		} else {
			return header, val, fmt.Errorf("For %s the value must either be set to Constant or Variable. The value is currently set to: %s", header, val)
		}
	}

	if header == "AIProcess" {
		if val == "Transcription" || val == "Object Recognition" || val == "Sports Classification" {
		} else {
			return header, val, fmt.Errorf("For %s the value must either be set to Transcription, Object Recognition or Sports Classification. The value is currently set to: %s", header, val)
		}
	}

	if header == "ContentCategories" {
		if val == "Demo Content" || val == "Case Studies" || val == "Promotional" || val == "Projects" || val == "Internal" || val == "Miscellaneous" {
		} else {
			return header, val, fmt.Errorf("For %s the value must either be set to Demo Content, Case Studies, Promotional, Projects, Internal or Miscellaneous. The value is currently set to: %s", header, val)
		}
	}

	if header == "win_ArchiveDelay" {
		_, err := strconv.Atoi(val)
		if err != nil {
			return header, val, fmt.Errorf("For %s the value must be set to an integer. The value is currently set to: %s", header, val)
		}
	}

	return header, val, nil
}

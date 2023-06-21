package reader

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
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
		// metadata["metadata_values"] = make(map[string]interface{})

		// Loop over each value in the row
		for count, value := range row {
			// First column is always our asset id
			if count == 0 {
				continue
				// TODO: Add validation for asset ID
			} else if count == 1 {
				// Second column is always our title
				title["title"] = value
			} else if count > 1 {

				// Columns after that are our metadata and variable in length
				header := fields[count]

				// Create a field value map
				fieldValue := map[string]interface{}{
					"value": value,
				}

				// Check if the header exists in metadataValues
				if _, ok := metadataValues[header]; !ok {
					// Create a new field values slice
					metadataValues[header] = map[string]interface{}{
						"field_values": []map[string]interface{}{},
					}
				}

				// Check if there is even anything in the column
				if value != "" {
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
		// Update the title
		updateTitle(cfg, row[0], title)

		// Assign the metadata values to the metadata map
		metadata["metadata_values"] = metadataValues

		// Update the metadata
		updateMetadata(cfg, row[0], metadata)
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

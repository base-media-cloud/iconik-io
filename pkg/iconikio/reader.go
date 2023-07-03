package iconikio

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
	"strings"

	"github.com/google/uuid"
)

// ReadCSVFile reads and validates the CSV file provided.
func (i *Iconik) ReadCSVFile() error {

	csvFile, err := os.Open(i.IconikClient.cfg.Input)
	if err != nil {
		return errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	fields, err := csvReader.Read()
	if err != nil {
		return errors.New("error reading CSV file")
	}

	if fields[0] != "id" || fields[1] != "title" {
		return errors.New("CSV file not properly formatted for Iconik")
	}

	csvColumnsName, _, err := i.GetCSVColumnsFromView()
	if err != nil {
		return err
	}
	csvColumnsName = append([]string{"title"}, csvColumnsName...)
	csvColumnsName = append([]string{"id"}, csvColumnsName...)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New("error reading CSV file")
		}

		title := make(map[string]string)
		metadata := make(map[string]interface{})
		metadataValues := make(map[string]interface{})

		for count, value := range row {
			if count == 0 {
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
				title["title"] = value
			} else if count > 1 {
				header := csvColumnsName[count]
				headerLabel := fields[count]

				if _, ok := metadataValues[header]; !ok {
					metadataValues[header] = map[string]interface{}{
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
							fieldValues := metadataValues[header].(map[string]interface{})["field_values"].([]map[string]interface{})
							fieldValues = append(fieldValues, fieldValue)
							metadataValues[header].(map[string]interface{})["field_values"] = fieldValues
						} else {
							delete(metadataValues, header)
						}
					}
				}
			}
		}
		err = updateTitle(i.IconikClient, row[0], title)
		if err != nil {
			return err
		}

		metadata["metadata_values"] = metadataValues

		err = updateMetadata(i.IconikClient, row[0], metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateTitle updates the title for the given asset ID.
func updateTitle(c *Client, assetID string, title map[string]string) error {

	requestBody, err := json.Marshal(title)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	uri := c.cfg.IconikURL + "/API/assets/v1/assets/" + assetID
	res, _, err := GetResponseBody("PATCH", uri, bytes.NewBuffer(requestBody), c)

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
func updateMetadata(c *Client, assetID string, metadata map[string]interface{}) error {

	requestBody, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	uri := c.cfg.IconikURL + "/API/metadata/v1/assets/" + assetID + "/views/" + c.cfg.ViewID + "/"
	res, _, err := GetResponseBody("PUT", uri, bytes.NewBuffer(requestBody), c)
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

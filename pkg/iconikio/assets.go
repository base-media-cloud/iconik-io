package iconikio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// GetCollectionAssets gets all the results from a collection and return the full object list with metadata.
func (i *Iconik) GetCollectionAssets() error {

	var a *Asset

	searchDoc := map[string]interface{}{
		"doc_types":        []string{"assets"},
		"query":            "",
		"metadata_view_id": i.IconikClient.Config.ViewID,
		"filter": map[string]interface{}{
			"operator": "AND",
			"terms": []map[string]interface{}{
				{
					"name":  "in_collections",
					"value": i.IconikClient.Config.CollectionID,
				},
			},
		},
	}

	requestBody, err := json.Marshal(searchDoc)
	if err != nil {
		return errors.New("error marshaling request body")
	}

	uri, err := i.joinURL("search", "", 0)
	if err != nil {
		return err
	}

	_, resBody, err := i.getResponseBody("POST", uri.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return err
	}

	dataNoNull := removeNullJSON(data)

	jsonData, err := json.MarshalIndent(dataNoNull, "", "  ")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &a)
	if err != nil {
		return err
	}

	i.IconikClient.Assets = append(i.IconikClient.Assets, a)

	return nil
}

// GetCSVColumnsFromView gets a column list from a metadata view for our CSV file, returning a slice of Names, and a slice of Labels.
func (i *Iconik) GetCSVColumnsFromView() ([]string, []string, error) {

	var csvColumnsName []string
	var csvColumnsLabel []string

	// uri := i.IconikClient.Config.IconikURL + "/API/metadata/v1/views/" + i.IconikClient.Config.ViewID

	uri, err := i.joinURL("metadataView", "", 0)
	if err != nil {
		return nil, nil, err
	}

	_, resBody, err := i.getResponseBody("GET", uri.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, nil, err
	}

	dataNoNull := removeNullJSON(data)

	jsonData, err := json.MarshalIndent(dataNoNull, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(jsonData, &i.IconikClient.Metadata)
	if err != nil {
		return nil, nil, err
	}

	for _, field := range i.IconikClient.Metadata.ViewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	return csvColumnsName, csvColumnsLabel, nil
}

func (i *Iconik) BuildCSVFile(csvColumnsName []string, csvColumnsLabel []string) error {
	// Get today's date and time
	today := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("%s.csv", today)
	filePath := i.IconikClient.Config.Output + filename

	// Open the CSV file
	csvFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("error creating CSV file")
	}
	defer csvFile.Close()

	metadataFile := csv.NewWriter(csvFile)
	defer metadataFile.Flush()

	// Write the header row
	headerRow := append([]string{"id", "title"}, csvColumnsLabel...)
	err = metadataFile.Write(headerRow)
	if err != nil {
		return errors.New("error writing header row")
	}
	numColumns := len(csvColumnsName)

	// Loop through all assets
	for _, asset := range i.IconikClient.Assets {
		for _, object := range asset.Objects {
			row := make([]string, numColumns+2) // +2 for id and title
			row[0] = object.ID
			row[1] = object.Title

			for i := 0; i < numColumns; i++ {
				metadataField := csvColumnsName[i]
				metadataValue, ok := object.Metadata[metadataField]
				if ok {
					switch v := metadataValue.(type) {
					case []interface{}:
						var values []string
						for _, val := range v {
							values = append(values, fmt.Sprintf("%v", val))
						}
						row[i+2] = strings.Join(values, ",")
					default:
						row[i+2] = fmt.Sprintf("%v", v)
					}
				}
			}

			err = metadataFile.Write(row)
			if err != nil {
				return errors.New("error writing row")
			}
		}
	}

	log.Println("File successfully saved to", filePath)
	return nil
}

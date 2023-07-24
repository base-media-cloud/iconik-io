package iconikio

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// GetCollection gets all the results from a collection and return the full object list with metadata.
func (i *Iconik) GetCollection() error {

	var c *Collection

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Path)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	_, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}

	jsonNoNull, err := removeNullJSONResBody(resBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonNoNull, &c)
	if err != nil {
		return err
	}

	if len(c.Errors) != 0 {
		return fmt.Errorf(strings.Join(c.Errors, ", "))
	}

	i.IconikClient.Collections = append(i.IconikClient.Collections, c)

	return nil
}

// GetMetadata gets the metadata using the given metadata view ID.
func (i *Iconik) GetMetadata() error {

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Path)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}

	err = IconikStatusCode(res)
	if err != nil {
		return err
	}

	jsonNoNull, err := removeNullJSONResBody(resBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonNoNull, &i.IconikClient.Metadata)
	if err != nil {
		return err
	}

	if len(i.IconikClient.Metadata.Errors) != 0 {
		return fmt.Errorf(strings.Join(i.IconikClient.Metadata.Errors, ", "))
	}

	return nil
}

func (i *Iconik) WriteCSVFile() error {

	// Get today's date and time
	today := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("%s.csv", today)
	filePath := i.IconikClient.Config.Output + filename

	// Create the CSV file
	csvFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("error creating CSV file")
	}
	defer csvFile.Close()

	metadataFile := csv.NewWriter(csvFile)
	defer metadataFile.Flush()

	var csvColumnsName []string
	var csvColumnsLabel []string

	for _, field := range i.IconikClient.Metadata.ViewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	// Write the header row
	headerRow := append([]string{"id", "title"}, csvColumnsLabel...)
	err = metadataFile.Write(headerRow)
	if err != nil {
		return errors.New("error writing header row")
	}
	numColumns := len(csvColumnsName)

	// Loop through all assets
	for _, collection := range i.IconikClient.Collections {
		for _, object := range collection.Objects {
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

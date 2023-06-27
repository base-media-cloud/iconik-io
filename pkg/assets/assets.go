package assets

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
	"time"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"go.uber.org/zap"
)

// get all results from a collection and return the full object list with metadata
func GetCollectionAssets(cfg *config.Conf, log *zap.SugaredLogger) (*Assets, error) {
	var assets *Assets
	url := cfg.IconikURL + "/API/search/v1/search/"
	log.Infow(url)
	method := "POST"

	searchDoc := map[string]interface{}{
		"doc_types":        []string{"assets"},
		"query":            "",
		"metadata_view_id": cfg.ViewID,
		"filter": map[string]interface{}{
			"operator": "AND",
			"terms": []map[string]interface{}{
				{
					"name":  "in_collections",
					"value": cfg.CollectionID,
				},
			},
		},
	}

	requestBody, err := json.Marshal(searchDoc)
	if err != nil {
		return nil, errors.New("error marshaling request body")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	queryParams := req.URL.Query()
	queryParams.Set("per_page", "150")
	queryParams.Set("scroll", "true")
	queryParams.Set("generate_signed_url", "false")
	queryParams.Set("save_search_history", "false")
	req.URL.RawQuery = queryParams.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

// get a column list from a metadata view for our CSV file
func GetCSVColumnsFromView(cfg *config.Conf, log *zap.SugaredLogger) ([]string, error) {

	var csvColumns []string
	var meta *MetadataFields

	url := cfg.IconikURL + "/API/metadata/v1/views/" + cfg.ViewID
	log.Infow(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &meta)
	if err != nil {
		return nil, err
	}

	for _, field := range meta.ViewFields {
		if field.Name != "__separator__" {
			csvColumns = append(csvColumns, field.Name)
		}
	}

	return csvColumns, nil
}

func (a *Assets) BuildCSVFile(cfg *config.Conf, csvColumns []string, log *zap.SugaredLogger) error {
	// Get today's date and time
	today := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("%s.csv", today)
	filePath := cfg.Output + filename

	// Open the CSV file
	csvFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("error creating CSV file")
	}
	defer csvFile.Close()

	metadataFile := csv.NewWriter(csvFile)
	defer metadataFile.Flush()

	// Write the header row
	headerRow := append([]string{"id", "title"}, csvColumns...)
	err = metadataFile.Write(headerRow)
	if err != nil {
		return errors.New("error writing header row")
	}
	numColumns := len(csvColumns)

	// Loop through all assets
	for _, asset := range a.Objects {
		row := make([]string, numColumns+2) // +2 for id and title
		row[0] = asset.ID
		row[1] = asset.Title

		for i := 0; i < numColumns; i++ {
			metadataField := csvColumns[i]
			metadataValue, ok := asset.Metadata[metadataField]
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

	log.Info("File successfully saved to", filePath)
	return nil
}

package assets

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/model"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/validate"
	"go.uber.org/zap"
)

// get all results from a collection and return the full object list with metadata
func GetCollectionAssets(cfg *config.Conf, log *zap.SugaredLogger) (*model.Assets, error) {
	var assets *model.Assets
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

	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, err
	}

	dataNoNull := validate.RemoveNullJSON(data)

	jsonData, err := json.MarshalIndent(dataNoNull, "", "  ")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

// get a column list from a metadata view for our CSV file
func GetCSVColumnsFromView(cfg *config.Conf, log *zap.SugaredLogger) ([]string, []string, error) {

	var csvColumnsName []string
	var csvColumnsLabel []string
	var meta *model.MetadataFields

	url := cfg.IconikURL + "/API/metadata/v1/views/" + cfg.ViewID
	log.Infow(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, nil, err
	}

	dataNoNull := validate.RemoveNullJSON(data)

	jsonData, err := json.MarshalIndent(dataNoNull, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(jsonData, &meta)
	if err != nil {
		return nil, nil, err
	}

	for _, field := range meta.ViewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	return csvColumnsName, csvColumnsLabel, nil
}

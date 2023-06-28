package assets

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/model"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/validate"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/web"
	"go.uber.org/zap"
)

// get all results from a collection and return the full object list with metadata
func GetCollectionAssets(cfg *config.Conf, log *zap.SugaredLogger) (*model.Assets, error) {
	var assets *model.Assets

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

	uri := cfg.IconikURL + "/API/search/v1/search/"
	_, resBody, err := web.GetResponseBody("POST", uri, bytes.NewBuffer(requestBody), cfg, log)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(resBody, &data)
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

	uri := cfg.IconikURL + "/API/metadata/v1/views/" + cfg.ViewID
	_, resBody, err := web.GetResponseBody("GET", uri, nil, cfg, log)
	if err != nil {
		return nil, nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(resBody, &data)
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

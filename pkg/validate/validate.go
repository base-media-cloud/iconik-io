package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/model"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/web"
	"go.uber.org/zap"
)

func CheckAppIDAuthTokenCollectionID(cfg *config.Conf, log *zap.SugaredLogger) error {

	uri := cfg.IconikURL + "/API/assets/v1/collections/" + cfg.CollectionID + "/contents/"
	res, _, err := web.GetResponseBody("GET", uri, nil, cfg, log)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized- please check your App-ID and Auth Token are correct")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("not found- please check your collection ID is correct")
	}

	return nil
}

func CheckMetadataID(cfg *config.Conf, log *zap.SugaredLogger) error {

	uri := cfg.IconikURL + "/API/metadata/v1/views/" + cfg.ViewID
	res, _, err := web.GetResponseBody("GET", uri, nil, cfg, log)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized- please check your metadata ID is correct")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("not found- please check your metadata ID is correct")
	}

	return nil
}

func CheckAssetbyID(assetID string, cfg *config.Conf, log *zap.SugaredLogger) (int, error) {

	uri := cfg.IconikURL + "/API/assets/v1/assets/" + assetID
	res, _, err := web.GetResponseBody("GET", uri, nil, cfg, log)
	if err != nil {
		return res.StatusCode, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return res.StatusCode, errors.New("unauthorized- please check your asset ID is correct")
	} else if res.StatusCode == http.StatusNotFound {
		return res.StatusCode, fmt.Errorf("%d: asset not found on Iconik servers", res.StatusCode)
	}

	return res.StatusCode, nil
}

func CheckAssetExistInCollection(assetID string, cfg *config.Conf, log *zap.SugaredLogger) (int, error) {
	var a *model.Assets
	uri := cfg.IconikURL + "/API/assets/v1/collections/" + cfg.CollectionID + "/contents/"
	res, resBody, err := web.GetResponseBody("GET", uri, nil, cfg, log)
	if err != nil {
		return res.StatusCode, err
	}

	err = json.Unmarshal(resBody, &a)
	if err != nil {
		return res.StatusCode, err
	}

	for _, asset := range a.Objects {
		if asset.ID == assetID {
			return res.StatusCode, nil
		}
	}

	return res.StatusCode, nil
}

func SchemaValidator(header, val string) (string, string, error) {

	if header == "Signed off" || header == "Archived" || header == "Can not be share" || header == "bmc_sapProductAssetOnly" {
		if val == "TRUE" {
			val = "true"
		} else if val == "FALSE" {
			val = "false"
		} else if val == "true" || val == "false" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to true or false. The value is currently set to: %s", header, val)
		}
	}

	if header == "Frame Rate" || header == "Audio Frame Rate" {
		if val == "23.976" || val == "23.98" || val == "24" || val == "25" || val == "29.97" || val == "30" || val == "50" || val == "59.94" || val == "60" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to 23.976, 23.98, 24, 25, 29.97, 30, 50, 59.94 or 60. The value is currently set to: %s", header, val)
		}
	}

	if header == "Frame Rate Mode" {
		if val == "Constant" || val == "Variable" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Constant or Variable. The value is currently set to: %s", header, val)
		}
	}

	if header == "AI Process" {
		if val == "Transcription" || val == "Object Recognition" || val == "Sports Classification" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Transcription, Object Recognition or Sports Classification. The value is currently set to: %s", header, val)
		}
	}

	if header == "Content Categories" {
		if val == "Demo Content" || val == "Case Studies" || val == "Promotional" || val == "Projects" || val == "Internal" || val == "Miscellaneous" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Demo Content, Case Studies, Promotional, Projects, Internal or Miscellaneous. The value is currently set to: %s", header, val)
		}
	}

	if header == "Archive Delay, days" {
		_, err := strconv.Atoi(val)
		if err != nil {
			return header, val, fmt.Errorf("for %s the value must be set to an integer. The value is currently set to: %s", header, val)
		}
	}

	return header, val, nil
}

func RemoveNullJSON(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v == nil {
			delete(m, k)
			continue
		}
		switch val := v.(type) {
		case map[string]interface{}:
			m[k] = RemoveNullJSON(val)
		case []interface{}:
			for i := 0; i < len(val); i++ {
				if _, ok := val[i].(map[string]interface{}); ok {
					val[i] = RemoveNullJSON(val[i].(map[string]interface{}))
				}
			}
		}
	}
	return m
}

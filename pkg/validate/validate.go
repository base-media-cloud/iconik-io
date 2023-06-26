package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/assets"
	"go.uber.org/zap"
)

func CheckAppIDAuthTokenCollectionID(cfg *config.Conf) error {
	// Check app ID, auth token and collection ID are all valid
	uri := cfg.IconikURL + "/API/assets/v1/collections/" + cfg.CollectionID + "/contents/?object_types=assets"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, nil)
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

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized- please check your App-ID and Auth Token are correct")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("not found- please check your collection ID is correct")
	}

	return nil
}

func CheckMetadataID(cfg *config.Conf) error {

	uri2 := cfg.IconikURL + "/API/metadata/v1/views/" + cfg.ViewID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri2, nil)
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

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized- please check your metadata ID is correct")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("not found- please check your metadata ID is correct")
	}

	return nil
}

func CheckAssetbyID(assetID string, cfg *config.Conf, log *zap.SugaredLogger) (int, error) {
	uri := cfg.IconikURL + "/API/assets/v1/assets/" + assetID
	log.Infow(uri)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return res.StatusCode, err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}

	if res.StatusCode == http.StatusNotFound {
		return res.StatusCode, fmt.Errorf("%d: Asset not found on Iconik servers", res.StatusCode)
	}

	return res.StatusCode, nil

}

func CheckAssetExistInCollection(assetID string, cfg *config.Conf, log *zap.SugaredLogger) (int, error) {
	var a *assets.Assets
	uri := cfg.IconikURL + "/API/assets/v1/collections/" + cfg.CollectionID + "/contents/?object_types=assets"
	log.Infow(uri)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}

	err = json.Unmarshal(responseBody, &a)
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
	// Schema validation for boolean fields
	if header == "Signedoff" || header == "win_Archived" || header == "ShareNo" || header == "bmc_sapProductAssetOnly" {
		if val == "TRUE" {
			val = "true"
		} else if val == "FALSE" {
			val = "false"
		} else if val == "true" || val == "false" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to true or false. The value is currently set to: %s", header, val)
		}
	}

	if header == "pdTest_FrameRate" || header == "pdTest_AudioFrameRate" {
		if val == "23.976" || val == "23.98" || val == "24" || val == "25" || val == "29.97" || val == "30" || val == "50" || val == "59.94" || val == "60" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to 23.976, 23.98, 24, 25, 29.97, 30, 50, 59.94 or 60. The value is currently set to: %s", header, val)
		}
	}

	if header == "pdTest_FrameRateMode" {
		if val == "Constant" || val == "Variable" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Constant or Variable. The value is currently set to: %s", header, val)
		}
	}

	if header == "AIProcess" {
		if val == "Transcription" || val == "Object Recognition" || val == "Sports Classification" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Transcription, Object Recognition or Sports Classification. The value is currently set to: %s", header, val)
		}
	}

	if header == "ContentCategories" {
		if val == "Demo Content" || val == "Case Studies" || val == "Promotional" || val == "Projects" || val == "Internal" || val == "Miscellaneous" {
		} else {
			return header, val, fmt.Errorf("for %s the value must either be set to Demo Content, Case Studies, Promotional, Projects, Internal or Miscellaneous. The value is currently set to: %s", header, val)
		}
	}

	if header == "win_ArchiveDelay" {
		_, err := strconv.Atoi(val)
		if err != nil {
			return header, val, fmt.Errorf("for %s the value must be set to an integer. The value is currently set to: %s", header, val)
		}
	}

	return header, val, nil
}

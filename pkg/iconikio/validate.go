package iconikio

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strconv"
)

// CheckAppIDAuthTokenCollectionID validates the App ID, Auth Token and Collection ID,
// and returns any errors to the user via the command line.
func (i *Iconik) CheckAppIDAuthTokenCollectionID() error {

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Path)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, _, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Method, u.String(), nil)
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

// CheckMetadataID validates the metadata view ID provided, and returns any errors to the user
// via the command line.
func (i *Iconik) CheckMetadataID() error {

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Path)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, _, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Method, u.String(), nil)
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

func (i *Iconik) validateAssetID(index int) error {
	// check asset id is valid
	_, err := uuid.Parse(i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	if err != nil {
		return errors.New("not a valid asset ID")
	}

	// check asset id exists on Iconik servers
	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Asset.Get.Path, i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	if err != nil {
		return err
	}
	u, err := url.Parse(result)
	if err != nil {
		return err
	}
	u.Scheme = i.IconikClient.Config.APIConfig.Scheme
	res, _, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Asset.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("ERROR: %d: UNAUTHORIZED, PLEASE CHECK YOUR ASSET ID %s IS CORRECT, SKIPPING", res.StatusCode, i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	} else if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("ERROR: %d: ASSET %s NOT FOUND ON ICONIK SERVERS, SKIPPING", res.StatusCode, i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	}

	// check asset id exists in given collection id
	var a *Asset
	result2, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Path)
	if err != nil {
		return err
	}
	u2, err := url.Parse(result2)
	if err != nil {
		return err
	}
	u2.Scheme = i.IconikClient.Config.APIConfig.Scheme
	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Method, u2.String(), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resBody, &a)
	if err != nil {
		return err
	}
	for _, object := range a.Objects {
		if object.ID == i.IconikClient.Config.CSVMetadata[index].IDStruct.ID {
			return nil
		}
	}
	return fmt.Errorf("ASSET %s DOES NOT EXIST IN GIVEN COLLECTION ID", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
}

// SchemaValidator checks values in the CSV against the matching headers, to see if they
// match the given criteria.
func SchemaValidator(header, val string) (string, string, error) {

	if header == "Signed off" || header == "Archived" || header == "Can not be share" || header == "SAP Product Asset" {
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

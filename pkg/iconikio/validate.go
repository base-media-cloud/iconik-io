package iconikio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// CheckAppIDAuthTokenCollectionID validates the App ID, Auth Token and Collection ID,
// and returns any errors to the user via the command line.
func (i *Iconik) CheckAppIDAuthTokenCollectionID() error {

	// uri := i.IconikClient.Config.IconikURL + "/API/assets/v1/collections/" + i.IconikClient.Config.CollectionID + "/contents/"

	uri, err := i.joinURL("collection", "", 0)
	if err != nil {
		return err
	}

	res, _, err := i.getResponseBody("GET", uri.String(), nil)
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

	// uri := i.IconikClient.Config.IconikURL + "/API/metadata/v1/views/" + i.IconikClient.Config.ViewID

	uri, err := i.joinURL("metadataView", "", 0)
	if err != nil {
		return err
	}

	res, _, err := i.getResponseBody("GET", uri.String(), nil)
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

// CheckAssetbyID validates the asset ID provided, and returns any errors to the user via
// the command line.
func (i *Iconik) CheckAssetbyID(assetID string) (int, error) {

	// uri := i.IconikClient.Config.IconikURL + "/API/assets/v1/assets/" + assetID

	uri, err := i.joinURL("asset", "", 0)
	if err != nil {
		return http.StatusNotFound, err
	}

	res, _, err := i.getResponseBody("GET", uri.String(), nil)
	if err != nil {
		return http.StatusNotFound, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return res.StatusCode, errors.New("unauthorized- please check your asset ID is correct")
	} else if res.StatusCode == http.StatusNotFound {
		return res.StatusCode, fmt.Errorf("%d: asset not found on Iconik servers", res.StatusCode)
	}

	return res.StatusCode, nil
}

// CheckAssetExistInCollection checks the asset ID against the collection ID provided, and
// returns any errors to the user via the command line.
func (i *Iconik) CheckAssetExistInCollection(assetID string) (int, error) {

	var a *Asset

	// uri := i.IconikClient.Config.IconikURL + "/API/assets/v1/collections/" + i.IconikClient.Config.CollectionID + "/contents/"
	uri, err := i.joinURL("collection", "", 0)
	if err != nil {
		return http.StatusNotFound, err
	}

	res, resBody, err := i.getResponseBody("GET", uri.String(), nil)
	if err != nil {
		return res.StatusCode, err
	}

	err = json.Unmarshal(resBody, &a)
	if err != nil {
		return res.StatusCode, err
	}

	for _, object := range a.Objects {
		if object.ID == assetID {
			return res.StatusCode, nil
		}
	}

	return res.StatusCode, nil
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

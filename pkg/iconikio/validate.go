package iconikio

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var validationRules = map[string][]string{
	"Signed off":          {"true", "false"},
	"Archived":            {"true", "false"},
	"Can not be share":    {"true", "false"},
	"SAP Product Asset":   {"true", "false"},
	"Frame Rate":          {"23.976", "23.98", "24", "25", "29.97", "30", "50", "59.94", "60"},
	"Audio Frame Rate":    {"23.976", "23.98", "24", "25", "29.97", "30", "50", "59.94", "60"},
	"Frame Rate Mode":     {"Constant", "Variable"},
	"AI Process":          {"Transcription", "Object Recognition", "Sports Classification"},
	"Content Categories":  {"Demo Content", "Case Studies", "Promotional", "Projects", "Internal", "Miscellaneous"},
	"Archive Delay, days": {}, // Empty slice means it should be an integer.
}

func SchemaValidator(header, val string) error {
	validValues, found := validationRules[header]
	if !found {
		return nil
	}

	if len(validValues) == 0 {
		_, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("for %s the value must be set to an integer. The value is currently set to: %s", header, val)
		}
		return nil
	}

	for _, validVal := range validValues {
		if val == validVal {
			return nil
		}
	}

	return fmt.Errorf("invalid value for %s. Valid values are: %s. The value is currently set to: %s", header, strings.Join(validValues, ", "), val)
}

func IconikStatusCode(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("status bad request")
	case http.StatusNotFound:
		return fmt.Errorf("status not found")
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorised- please check your App ID and Auth Token are correct")
	default:
		return nil
	}
}

func (i *Iconik) validateFilename(index int) error {
	// check filename exists in given collection id
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
	err = json.Unmarshal(resBody, &c)
	if err != nil {
		return err
	}
	for _, object := range c.Objects {
		for _, file := range object.Files {
			if file.OriginalName == i.IconikClient.Config.CSVMetadata[index].OriginalNameStruct.OriginalName {
				i.IconikClient.Config.CSVMetadata[index].IDStruct.ID = object.ID
				return nil
			}
		}
	}
	return fmt.Errorf("file %s does not exist in given collection id", i.IconikClient.Config.CSVMetadata[index].OriginalNameStruct.OriginalName)
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
	var c *Collection
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
	err = json.Unmarshal(resBody, &c)
	if err != nil {
		return err
	}
	for _, object := range c.Objects {
		if object.ID == i.IconikClient.Config.CSVMetadata[index].IDStruct.ID {
			return nil
		}
	}
	return fmt.Errorf("ASSET %s DOES NOT EXIST IN GIVEN COLLECTION ID", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
}

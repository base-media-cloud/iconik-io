package validate

import (
	"context"
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	csvdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/csv"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

// API is an interface that defines the operations that can be performed on the system domains endpoint.
type API interface {
	GetAsset(ctx context.Context, path, assetID string) (*assets.DTO, error)
}

// Svc is a struct that implements the systemdomainports.Servicer interface.
type Svc struct {
	api API
}

// New is a function that returns a new instance of the Validator struct.
func New(api API) *Svc {
	return &Svc{
		api: api,
	}
}

func (svc *Svc) Schema(header, val string) error {
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

func (svc *Svc) Filename(objects []collections.Object, csvMetadata csvdomain.CSVMetadata) error {
	for _, object := range objects {
		for _, file := range object.Files {
			if file.OriginalName == csvMetadata.OriginalNameStruct.OriginalName {
				csvMetadata.IDStruct.ID = object.ID
				return nil
			}
		}
	}

	return fmt.Errorf("file %s does not exist in given collection id", csvMetadata.OriginalNameStruct.OriginalName)
}

func (svc *Svc) AssetID(objects []collections.Object, csvMetadata csvdomain.CSVMetadata, ctx context.Context) error {
	_, err := uuid.Parse(csvMetadata.IDStruct.ID)
	if err != nil {
		return errors.New("not a valid asset ID")
	}

	_, err = svc.api.GetAsset(ctx, iconik.AssetsPath, csvMetadata.IDStruct.ID)
	if err != nil {
		return fmt.Errorf("asset %s not found on iconik servers", csvMetadata.IDStruct.ID)
	}

	for _, object := range objects {
		for _, file := range object.Files {
			if object.ID == csvMetadata.IDStruct.ID {
				csvMetadata.OriginalNameStruct.OriginalName = file.OriginalName
				return nil
			}
		}
	}

	return fmt.Errorf("asset %s does not exist in given collection id", csvMetadata.IDStruct.ID)
}

// func IconikStatusCode(res *http.Response) error {
// 	switch res.StatusCode {
// 	case http.StatusBadRequest:
// 		return fmt.Errorf("status bad request")
// 	case http.StatusNotFound:
// 		return fmt.Errorf("status not found")
// 	case http.StatusUnauthorized:
// 		return fmt.Errorf("unauthorised- please check your App ID and Auth Token are correct")
// 	default:
// 		return nil
// 	}
// }

package utils

import (
	"fmt"
	csvdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/csv"
	collDomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
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

func ValidateSchema(header, val string) error {
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

func ValidateFilename(objects []collDomain.ObjectDTO, csvMetadata csvdomain.CSVMetadata) error {
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

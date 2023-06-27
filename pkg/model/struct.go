package model

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"go.uber.org/zap"
)

// ====================================================
// iconik Objects Response Structure "GET /v1/assets/"

// Assets is the top level data structure that receives the unmarshalled payload
// response.
type Assets struct {
	Objects []*Object `json:"objects"`
}

// Objects acts as a non nested struct to the Objects type in Assets.
type Object struct {
	ID       string                 `json:"id"`
	Metadata map[string]interface{} `json:"metadata"`
	Title    string                 `json:"title"`
}

// ====================================================
// iconik Objects Response Structure "GET /API/metadata/v1/views/"

// MetadataFields is the top level data structure that receives the unmarshalled payload
// response.
type MetadataFields struct {
	ViewFields []*ViewField `json:"view_fields"`
}

// ViewField acts as a non nested struct to the ViewFields type in MetadataFields.
type ViewField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

func (a *Assets) BuildCSVFile(cfg *config.Conf, csvColumnsName []string, csvColumnsLabel []string, log *zap.SugaredLogger) error {
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
	headerRow := append([]string{"id", "title"}, csvColumnsLabel...)
	err = metadataFile.Write(headerRow)
	if err != nil {
		return errors.New("error writing header row")
	}
	numColumns := len(csvColumnsName)

	// Loop through all assets
	for _, asset := range a.Objects {
		row := make([]string, numColumns+2) // +2 for id and title
		row[0] = asset.ID
		row[1] = asset.Title

		for i := 0; i < numColumns; i++ {
			metadataField := csvColumnsName[i]
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

package iconikio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/url"
	"os"
	"strings"
)

func (i *Iconik) ReadCSVFile() ([][]string, error) {
	csvFile, err := os.Open(i.IconikClient.Config.Input)
	if err != nil {
		return nil, errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData, nil
}

func (i *Iconik) ReadExcelFile() ([][]string, error) {
	excelFile, err := excelize.OpenFile(i.IconikClient.Config.Input)
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()

	activeSheet := excelFile.GetSheetList()[excelFile.GetActiveSheetIndex()]

	excelData, err := excelFile.GetRows(activeSheet)
	if err != nil {
		return nil, err
	}

	return excelData, nil
}

// UpdateIconik reads a 2D slice, verifies it, and uploads the data to the Iconik API.
func (i *Iconik) UpdateIconik(metadataFile [][]string) error {
	csvHeaders := metadataFile[0]

	if csvHeaders[0] != "id" || csvHeaders[1] != "original_name" || csvHeaders[2] != "title" {
		fmt.Println(csvHeaders)
		return errors.New("CSV file not properly formatted for Iconik")
	}

	matchingData, nonMatchingHeaders, err := i.matchCSVtoAPI(metadataFile)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		fmt.Println("Some columns from the file provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
		fmt.Println()
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	matchingFileHeaderNames := matchingData[0]
	matchingFileHeaderLabels := matchingData[1]

	i.IconikClient.Config.CSVFilesToUpdate = len(matchingData) - 2
	fmt.Println("Amount of files to update:", i.IconikClient.Config.CSVFilesToUpdate)

	for index := 2; index < len(matchingData); index++ {
		row := matchingData[index]

		csvMetadata := CSVMetadata{
			Added: false,
			IDStruct: IDStruct{
				ID: row[0],
			},
			OriginalNameStruct: OriginalNameStruct{
				OriginalName: row[1],
			},
			TitleStruct: TitleStruct{
				Title: row[2],
			},
			MetadataValuesStruct: MetadataValuesStruct{
				MetadataValues: make(map[string]struct {
					FieldValues []FieldValue `json:"field_values"`
				}),
			},
		}

		i.IconikClient.Config.CSVMetadata = append(i.IconikClient.Config.CSVMetadata, &csvMetadata)

		err := i.validateAssetID(index - 2)
		err2 := i.validateFilename(index - 2)
		if err != nil && err2 != nil {
			log.Printf("%s & %s, skipping\n", err, err2)
			continue
		}
		csvMetadata.Added = true

		for count := 3; count < len(row); count++ {
			headerName := matchingFileHeaderNames[count]
			headerLabel := matchingFileHeaderLabels[count]
			fieldValueSlice := make([]FieldValue, 0)

			valueArr := strings.Split(row[count], ",")
			if !isBlankStringArray(valueArr) {
				for _, val := range valueArr {

					err = SchemaValidator(headerLabel, val)
					if err != nil {
						return err
					}

					fieldValueSlice = append(fieldValueSlice, FieldValue{Value: val})
				}
				csvMetadata.MetadataValuesStruct.MetadataValues[headerName] = struct {
					FieldValues []FieldValue `json:"field_values"`
				}{
					FieldValues: fieldValueSlice,
				}
			} else {
				continue
			}
		}

		err = i.updateTitle(index - 2)
		if err != nil {
			return err
		}

		err = i.updateMetadata(index - 2)
		if err != nil {
			return err
		}

	}

	fmt.Println()
	log.Println("Assets successfully updated:")
	var countSuccess int
	for _, csvMetadata := range i.IconikClient.Config.CSVMetadata {
		if csvMetadata.Added {
			countSuccess++
		}
	}
	fmt.Printf("%d of %d", countSuccess, i.IconikClient.Config.CSVFilesToUpdate)

	fmt.Println()
	log.Println("Assets that failed to update:")
	var countFailed int
	for _, csvMetadata := range i.IconikClient.Config.CSVMetadata {
		if !csvMetadata.Added {
			countFailed++
			log.Printf("%s (Title: %s, Original filename: %s)", csvMetadata.IDStruct.ID, csvMetadata.TitleStruct.Title, csvMetadata.OriginalNameStruct.OriginalName)
		}
	}
	fmt.Printf("%d of %d\n", countFailed, i.IconikClient.Config.CSVFilesToUpdate)

	return nil
}

func (i *Iconik) updateTitle(index int) error {
	requestBody, err := json.Marshal(i.IconikClient.Config.CSVMetadata[index].TitleStruct)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Asset.Patch.Path, i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	if err != nil {
		return err
	}
	u, err := url.Parse(result)
	if err != nil {
		return err
	}
	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Asset.Patch.Method, u.String(), bytes.NewBuffer(requestBody))

	if res.StatusCode == 200 {
	} else {
		log.Println("Error updating title name for asset ", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
		log.Println("resBody:", string(resBody))
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

func (i *Iconik) updateMetadata(index int) error {
	requestBody, err := json.Marshal(i.IconikClient.Config.CSVMetadata[index].MetadataValuesStruct)
	if err != nil {
		return errors.New("error marshaling JSON")
	}

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Path, i.IconikClient.Config.CSVMetadata[index].IDStruct.ID, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Path2)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Put.Method, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		log.Printf("Successfully updated metadata for asset %s (%s)", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID, i.IconikClient.Config.CSVMetadata[index].OriginalNameStruct.OriginalName)
	} else {
		log.Println("Error updating metadata for asset ", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
		log.Println("resBody:", string(resBody))
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

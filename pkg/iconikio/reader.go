package iconikio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

// ReadCSVFile reads and validates the CSV file provided.
func (i *Iconik) ReadCSVFile() error {

	// open the provided CSV file
	csvFile, err := os.Open(i.IconikClient.Config.Input)
	if err != nil {
		return errors.New("error opening CSV file")
	}
	defer csvFile.Close()

	// create a new reader
	csvReader := csv.NewReader(csvFile)

	// read all the CSV data into a 2D slice we can work with
	csvData, err := csvReader.ReadAll()
	if err != nil {
		//return errors.New("error reading CSV file")
		return err
	}

	// the first row of the 2D slice will be the header row
	csvHeaders := csvData[0]

	// we then validate that the schema for the header row is correct
	if csvHeaders[0] != "id" || csvHeaders[1] != "original_name" || csvHeaders[2] != "title" {
		return errors.New("CSV file not properly formatted for Iconik")
	}

	// get the slimmed down 2D slice that contains only the columns matched to the given metadata view
	matchingCSV, nonMatchingHeaders, err := i.matchCSVtoAPI(csvData)
	if err != nil {
		return err
	}

	if len(nonMatchingHeaders) > 0 {
		// log to the user if there were any columns in the provided CSV that have been left out
		fmt.Println("Some columns from the CSV provided have not been included in the upload to Iconik, as they are not part of the metadata view provided. Please see below for the headers of the columns not included:")
		fmt.Println()
		for _, nonMatchingHeader := range nonMatchingHeaders {
			fmt.Println(nonMatchingHeader)
		}
	}

	// get the CSV Header Names row
	matchingCSVHeaderNames := matchingCSV[0]
	// get the CSV Header Labels row
	matchingCSVHeaderLabels := matchingCSV[1]

	for index := 2; index < len(matchingCSV); index++ {
		fmt.Println()
		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
		row := matchingCSV[index]

		// Create a new instance of CSVMetadata.
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

		log.Printf("Attempting to update metadata for asset ID %s from row %d of the provided CSV:", csvMetadata.IDStruct.ID, index-1)

		err := i.validateAssetID(index - 2)
		err2 := i.validateFilename(index - 2)
		if err != nil && err2 != nil {
			log.Printf("%s & %s, skipping\n", err, err2)
			fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
			continue
		}
		csvMetadata.Added = true

		for count := 3; count < len(row); count++ {
			headerName := matchingCSVHeaderNames[count]
			headerLabel := matchingCSVHeaderLabels[count]
			fieldValueSlice := make([]FieldValue, 0)

			// if there are words separated by spaces inside the field, split with a comma
			valueArr := strings.Split(row[count], ",")
			if !isBlankStringArray(valueArr) {
				for _, val := range valueArr {

					err = SchemaValidator(headerLabel, val)
					if err != nil {
						return err
					}

					fieldValueSlice = append(fieldValueSlice, FieldValue{Value: val})
				}
				// Add the fieldValueSlice to the specific headerName in MetadataValues
				csvMetadata.MetadataValuesStruct.MetadataValues[headerName] = struct {
					FieldValues []FieldValue `json:"field_values"`
				}{
					FieldValues: fieldValueSlice,
				}
			} else {
				continue
			}
		}

		// update the title name for this particular asset/row
		err = i.updateTitle(index - 2)
		if err != nil {
			return err
		}

		// update the remaining metadata for this particular asset/row
		err = i.updateMetadata(index - 2)
		if err != nil {
			return err
		}

		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
	}

	fmt.Println()
	log.Println("Assets successfully updated:")
	for _, csvMetadata := range i.IconikClient.Config.CSVMetadata {
		if csvMetadata.Added {
			log.Printf("%s (Title: %s, Original filename: %s)", csvMetadata.IDStruct.ID, csvMetadata.TitleStruct.Title, csvMetadata.OriginalNameStruct.OriginalName)
		}
	}

	fmt.Println()
	log.Println("Assets that failed to update:")
	for _, csvMetadata := range i.IconikClient.Config.CSVMetadata {
		if !csvMetadata.Added {
			log.Printf("%s (Title: %s, Original filename: %s)", csvMetadata.IDStruct.ID, csvMetadata.TitleStruct.Title, csvMetadata.OriginalNameStruct.OriginalName)
		}
	}
	fmt.Println()

	return nil
}

// updateTitle updates the title for the given asset ID.
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
		log.Printf("Successfully updated title name for asset %s (%s)", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID, i.IconikClient.Config.CSVMetadata[index].OriginalNameStruct.OriginalName)
	} else {
		log.Println("Error updating title name for asset ", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
		log.Println("resBody:", string(resBody))
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

// updateMetadata updates the metadata for the given asset ID.
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
		log.Println("Successfully updated metadata for asset ", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
	} else {
		log.Println("Error updating metadata for asset ", i.IconikClient.Config.CSVMetadata[index].IDStruct.ID)
		log.Println("resBody:", string(resBody))
		log.Println(fmt.Sprint(res.StatusCode))
		return err
	}

	return nil
}

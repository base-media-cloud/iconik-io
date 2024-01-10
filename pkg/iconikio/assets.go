package iconikio

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// GetCollection gets all the results from a collection and return the full object list with metadata.
func (i *Iconik) GetCol(collectionID string, pageNo int, w *csv.Writer) error {
	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Path, collectionID, "/contents/")
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme
	queryParams := u.Query()
	queryParams.Set("per_page", "500")
	queryParams.Set("page", strconv.Itoa(pageNo))
	u.RawQuery = queryParams.Encode()

	_, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.Collection.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}

	var c *Collection
	if err = json.Unmarshal(resBody, &c); err != nil {
		return err
	}

	if err = i.ProcessObjs(c, w); err != nil {
		return err
	}

	if c.Pages > pageNo {
		if err = i.GetCol(collectionID, pageNo+1, w); err != nil {
			return err
		}
	}

	return nil
}

func (i *Iconik) ProcessObjs(c *Collection, w *csv.Writer) error {
	var output []*Object

	for j := range c.Objects {
		if c.Objects[j].ObjectType == "collections" {
			fmt.Printf("\nfound collection %s", c.Objects[j].Title)
			if err := i.GetCol(c.Objects[j].ID, 1, w); err != nil {
				return err
			}
			if err := i.ProcessObjs(c, w); err != nil {
				return err
			}
			continue
		}
		output = append(output, c.Objects[j])
	}

	toWrite, err := i.PrepMetadata(output)
	if err != nil {
		return err
	}

	if err = w.WriteAll(toWrite); err != nil {
		return err
	}

	return nil
}

// GetMetadata gets the metadata using the given metadata view ID.
func (i *Iconik) GetMetadata() error {
	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Path)
	if err != nil {
		return err
	}

	u, err := url.Parse(result)
	if err != nil {
		return err
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	res, resBody, err := i.getResponseBody(i.IconikClient.Config.APIConfig.Endpoints.MetadataView.Get.Method, u.String(), nil)
	if err != nil {
		return err
	}

	err = IconikStatusCode(res)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, &i.IconikClient.Metadata)
	if err != nil {
		return err
	}

	if len(i.IconikClient.Metadata.Errors) != 0 {
		return fmt.Errorf(strings.Join(i.IconikClient.Metadata.Errors, ", "))
	}

	return nil
}

func (i *Iconik) PrepMetadataForWriting() ([][]string, error) {
	var metadataFile [][]string
	var csvColumnsName []string
	var csvColumnsLabel []string

	for _, field := range i.IconikClient.Metadata.ViewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	// Write the header row
	headerRow := append([]string{"id", "original_name", "size", "title"}, csvColumnsLabel...)
	metadataFile = append(metadataFile, headerRow)
	numColumns := len(csvColumnsName)

	// Loop through all assets
	for _, object := range i.IconikClient.Assets {
		row := make([]string, numColumns+4)
		row[0] = object.ID
		row[1] = object.Files[0].OriginalName
		row[2] = strconv.Itoa(object.Files[0].Size)
		row[3] = object.Title

		for i := 0; i < numColumns; i++ {
			metadataField := csvColumnsName[i]
			metadataValue := object.Metadata[metadataField]
			result := make([]string, len(metadataValue))

			for index, elem := range metadataValue {
				switch val := elem.(type) {
				case string:
					str := val
					if strings.HasPrefix(str, " ") {
						str = strings.TrimLeft(str, " ")
					}
					if strings.HasSuffix(str, " ") {
						str = strings.TrimRight(str, " ")
					}
					result[index] = str
				case bool:
					result[index] = fmt.Sprintf("%t", val)
				case int:
					result[index] = fmt.Sprintf("%d", val)
				default:
				}
			}

			if len(result) > 1 {
				row[i+4] = strings.Join(result, ",")
			} else {
				row[i+4] = strings.Join(result, "")
			}

		}

		metadataFile = append(metadataFile, row)
	}

	return metadataFile, nil
}

func (i *Iconik) PrepMetadata(objs []*Object) ([][]string, error) {
	var metadataFile [][]string
	var csvColumnsName []string
	var csvColumnsLabel []string

	for _, field := range i.IconikClient.Metadata.ViewFields {
		if field.Name != "__separator__" {
			csvColumnsName = append(csvColumnsName, field.Name)
			csvColumnsLabel = append(csvColumnsLabel, field.Label)
		}
	}

	// Write the header row
	headerRow := append([]string{"id", "original_name", "size", "title"}, csvColumnsLabel...)
	metadataFile = append(metadataFile, headerRow)
	numColumns := len(csvColumnsName)

	// Loop through all assets
	for _, object := range objs {
		row := make([]string, numColumns+4)
		row[0] = object.ID
		row[1] = object.Files[0].OriginalName
		row[2] = strconv.Itoa(object.Files[0].Size)
		row[3] = object.Title

		for i := 0; i < numColumns; i++ {
			metadataField := csvColumnsName[i]
			metadataValue := object.Metadata[metadataField]
			result := make([]string, len(metadataValue))

			for index, elem := range metadataValue {
				switch val := elem.(type) {
				case string:
					str := val
					if strings.HasPrefix(str, " ") {
						str = strings.TrimLeft(str, " ")
					}
					if strings.HasSuffix(str, " ") {
						str = strings.TrimRight(str, " ")
					}
					result[index] = str
				case bool:
					result[index] = fmt.Sprintf("%t", val)
				case int:
					result[index] = fmt.Sprintf("%d", val)
				default:
				}
			}

			if len(result) > 1 {
				row[i+4] = strings.Join(result, ",")
			} else {
				row[i+4] = strings.Join(result, "")
			}

		}

		metadataFile = append(metadataFile, row)
	}

	return metadataFile, nil
}

func (i *Iconik) WriteCSVFile(metadataFile [][]string) error {
	today := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("%s.csv", today)
	filePath := i.IconikClient.Config.Output + filename

	// Create the CSV file
	csvFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("error creating CSV file")
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	err = csvWriter.WriteAll(metadataFile)
	if err != nil {
		return err
	}

	log.Println("CSV file successfully saved to", filePath)
	return nil
}

func (i *Iconik) WriteExcelFile(metadataFile [][]string) error {
	today := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("%s.xlsx", today)
	filePath := i.IconikClient.Config.Output + filename
	sheetName := today

	// Create the excel file
	excelFile := excelize.NewFile()
	defer excelFile.Close()
	if err := excelFile.SetSheetName("Sheet1", sheetName); err != nil {
		return err
	}

	for i, row := range metadataFile {
		startCell, err := excelize.JoinCellName("A", i+1)
		if err != nil {
			return err
		}
		if err := excelFile.SetSheetRow(sheetName, startCell, &row); err != nil {
			return err
		}
	}

	if err := excelFile.SaveAs(filePath); err != nil {
		return err
	}

	log.Println("Excel file successfully saved to", filePath)
	return nil
}

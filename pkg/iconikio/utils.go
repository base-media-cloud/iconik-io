package iconikio

import (
	"io"
	"net/http"
)

func (i *Iconik) matchCSVtoAPI(csvData [][]string) ([][]string, []string, error) {

	csvHeaderLabels := csvData[0]

	var matchingIconikHeaderNames []string
	var matchingIconikHeaderLabels []string
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "id")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "original_name")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "title")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "id")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "original_name")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "title")

	var nonMatchingHeaders []string

	for index, csvHeaderLabel := range csvHeaderLabels {
		if index > 2 {
			found := false
			for _, field := range i.IconikClient.Metadata.ViewFields {
				if csvHeaderLabel == field.Label {
					matchingIconikHeaderNames = append(matchingIconikHeaderNames, field.Name)
					matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, field.Label)
					found = true
					break
				}
			}
			if !found {
				nonMatchingHeaders = append(nonMatchingHeaders, csvHeaderLabel)
			}
		}
	}

	var matchingValues [][]string
	matchingValues = append(matchingValues, matchingIconikHeaderNames)
	matchingValues = append(matchingValues, matchingIconikHeaderLabels)

	for i := 1; i < len(csvData); i++ {
		row := csvData[i]
		var matchingRow []string
		for i, csvHeaderLabel := range csvHeaderLabels {
			if contains(matchingIconikHeaderLabels, csvHeaderLabel) {
				matchingRow = append(matchingRow, row[i])
			}
		}
		matchingValues = append(matchingValues, matchingRow)
	}

	return matchingValues, nonMatchingHeaders, nil

}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func isBlankStringArray(arr []string) bool {
	for _, s := range arr {
		if s != "" {
			return false
		}
	}
	return true
}

func (i *Iconik) getResponseBody(method, uri string, params io.Reader) (*http.Response, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, uri, params)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("App-ID", i.IconikClient.Config.AppID)
	req.Header.Add("Auth-Token", i.IconikClient.Config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return res, resBody, nil
}

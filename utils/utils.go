package utils

import (
	metadataDomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
)

func MatchCSVtoAPI(viewFields []metadataDomain.ViewFieldDTO, csvData [][]string) ([][]string, []string, error) {

	csvHeaderLabels := csvData[0]

	var matchingIconikHeaderNames []string
	var matchingIconikHeaderLabels []string
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "id")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "original_name")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "size")
	matchingIconikHeaderNames = append(matchingIconikHeaderNames, "title")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "id")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "original_name")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "size")
	matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, "title")

	var nonMatchingHeaders []string

	for index, csvHeaderLabel := range csvHeaderLabels {
		if index > 3 {
			found := false
			for _, viewField := range viewFields {
				if csvHeaderLabel == viewField.Label {
					matchingIconikHeaderNames = append(matchingIconikHeaderNames, viewField.Name)
					matchingIconikHeaderLabels = append(matchingIconikHeaderLabels, viewField.Label)
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

	for j := 1; j < len(csvData); j++ {
		row := csvData[j]
		var matchingRow []string
		for k, csvHeaderLabel := range csvHeaderLabels {
			if contains(matchingIconikHeaderLabels, csvHeaderLabel) {
				matchingRow = append(matchingRow, row[k])
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

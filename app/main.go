package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/reader"
)

type CMDArgs struct {
	IconikURL    string
	AppID        string
	AuthToken    string
	CollectionID string
	ViewID       string
	Input        string
	Output       string
}

var (
	cfg  config.Conf
	cmds CMDArgs
)

func main() {
	// Parse the flags entered
	argParse()

	// Construct config struct from command line args
	constructConfig(&cmds)

	if cmds.Output != "" {
		// User has chosen CSV output

		// Get Assets
		a, err := assets.GetCollectionAssets(&cfg)
		if err != nil {
			panic(err)
		}

		// Get CSV Headers
		columns, err := assets.GetCSVColumnsFromView(&cfg)
		if err != nil {
			panic(err)
		}

		// Build CSV and output
		err = a.BuildCSVFile(&cfg, columns)
		if err != nil {
			panic(err)
		}
	}

	if cmds.Input != "" {
		// User has chosen CSV input
		err := reader.ReadCSVFile(&cfg)
		if err != nil {
			panic(err)
		}
	}
}

func argParse() {
	flag.StringVar(&cmds.IconikURL, "iconik-url", "https://preview.iconik.cloud", "iconik URL")
	flag.StringVar(&cmds.AppID, "app-id", "", "iconik Application ID")
	flag.StringVar(&cmds.AuthToken, "auth-token", "", "iconik Authentication token")
	flag.StringVar(&cmds.CollectionID, "collection-id", "", "iconik Collection ID")
	flag.StringVar(&cmds.ViewID, "metadata-view-id", "", "iconik Metadata View ID")
	flag.StringVar(&cmds.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cmds.Output, "output", "", "Output mode - requires path to save CSV file")
	flag.Parse()

	if cmds.AppID == "" {
		log.Fatal("No App-Id provided")
	}
	if cmds.AuthToken == "" {
		log.Fatal("No Auth-Token provided")
	}
	if cmds.CollectionID == "" {
		log.Fatal("No Collection ID provided")
	}
	if cmds.ViewID == "" {
		log.Fatal("No Metadata View ID provided")
	}
	if cmds.Input == "" && cmds.Output == "" {
		log.Fatal("Neither input or output mode selected. Please select one.")
	}

	err := iconikErrorCheck(&cmds)
	if err != nil {
		panic(err)
	}

	err = iconikErrorCheckMetadata(&cmds)
	if err != nil {
		panic(err)
	}

}

func constructConfig(args *CMDArgs) {
	cfg.IconikURL = args.IconikURL
	cfg.AppID = args.AppID
	cfg.AuthToken = args.AuthToken
	cfg.CollectionID = args.CollectionID
	cfg.ViewID = args.ViewID
	cfg.Input = args.Input
	cfg.Output = args.Output
}

func iconikErrorCheck(args *CMDArgs) error {
	// Check app ID, auth token and collection ID are all valid
	uri := cmds.IconikURL + "/API/assets/v1/collections/" + cmds.CollectionID + "/contents/?object_types=assets"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("App-ID", cmds.AppID)
	req.Header.Add("Auth-Token", cmds.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.New("Unauthorized. Please check your App-ID and Auth Token are correct.")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("Not found. Please check your collection ID is correct.")
	}

	return nil
}

func iconikErrorCheckMetadata(args *CMDArgs) error {

	uri2 := cmds.IconikURL + "/API/metadata/v1/views/" + cmds.ViewID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri2, nil)
	if err != nil {
		return err
	}

	req.Header.Add("App-ID", cmds.AppID)
	req.Header.Add("Auth-Token", cmds.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.New("Unauthorized. Please check your metadata ID is correct.")
	} else if res.StatusCode == http.StatusNotFound {
		return errors.New("Not found. Please check your metadata ID is correct.")
	}

	return nil
}

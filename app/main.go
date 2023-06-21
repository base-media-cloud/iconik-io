package main

import (
	"flag"
	"fmt"
	"log"

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
		fmt.Println(a)

		// Get CSV Headers
		columns, err := assets.GetCSVColumnsFromView(&cfg)
		if err != nil {
			panic(err)
		}

		// Build CSV and output
		a.BuildCSVFile(&cfg, columns)
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
	flag.StringVar(&cmds.AppID, "app-id", "", "The iconik Application ID")
	flag.StringVar(&cmds.AuthToken, "auth-token", "", "The iconik Authentication token")
	flag.StringVar(&cmds.CollectionID, "collection-id", "", "iconik Collection Id")
	flag.StringVar(&cmds.ViewID, "metadata-view-id", "", "iconik Metadata View Id")
	flag.StringVar(&cmds.Input, "input", "", "Input mode - requires path to input CSV file")
	flag.StringVar(&cmds.Output, "output", "csv/", "Output mode - requires path to save CSV file")
	flag.Parse()
	if cmds.AppID == "" {
		log.Fatal("No App-Id provided")
	}
	if cmds.AuthToken == "" {
		log.Fatal("No Auth-Token provided")
	}
	if cmds.IconikURL == "" {
		log.Fatal("No Iconik URL provided")
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

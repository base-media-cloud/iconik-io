package cmd

import (
	"flag"
	"log"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/assets"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/reader"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/validate"
	"go.uber.org/zap"
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

func Execute(log *zap.SugaredLogger) error {

	// Parse the flags entered
	err := argParse()
	if err != nil {
		return err
	}

	if cfg.Output != "" {
		// User has chosen CSV output

		// Get Assets
		a, err := assets.GetCollectionAssets(&cfg, log)
		if err != nil {
			return err
		}

		// Get CSV Headers
		columns, err := assets.GetCSVColumnsFromView(&cfg, log)
		if err != nil {
			return err
		}

		// Build CSV and output
		err = a.BuildCSVFile(&cfg, columns, log)
		if err != nil {
			return err
		}
	}

	if cfg.Input != "" {
		// User has chosen CSV input
		err := reader.ReadCSVFile(&cfg, log)
		if err != nil {
			return err
		}
	}

	return nil
}

func argParse() error {
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

	// Construct config struct from command line args
	constructConfig(&cmds)

	err := validate.CheckAppIDAuthTokenCollectionID(&cfg)
	if err != nil {
		return err
	}

	err = validate.CheckMetadataID(&cfg)
	if err != nil {
		return err
	}

	return nil

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

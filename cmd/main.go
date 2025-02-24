package main

import (
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/app/input"
	"github.com/base-media-cloud/pd-iconik-io-rd/app/output"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/api/iconik"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	assetsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/iconik/assets/assets"
	collsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/iconik/assets/collections"
	metadatasvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/iconik/metadata"
	searchsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/iconik/search"
	inputsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/input"
	outputsvc "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/services/output"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logFilenameTimeStamped := fmt.Sprintf("iconik-io_%s.log", time.Now().Format(time.DateOnly))
	logFile, err := os.OpenFile(logFilenameTimeStamped, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error creating log file")
		log.Fatal(err)
	}
	defer logFile.Close()

	l := logger.New(logFile)

	cfg, err := config.NewApp()
	if err != nil {
		fmt.Println(domain.ErrInternalError)
		l.Fatal().Err(err).Msg("error creating app config")
	}

	req := api.New(&http.Client{})
	iconikAPI := iconik.New(cfg, req)

	assetSvc := assetsvc.New(iconikAPI)
	collSvc := collsvc.New(iconikAPI)
	metadataSvc := metadatasvc.New(iconikAPI)
	searchSvc := searchsvc.New(iconikAPI)

	if cfg.Type == input.AppType {
		inputSvc := inputsvc.New(collSvc, assetSvc, metadataSvc, searchSvc)
		if err = input.Run(cfg, inputSvc, l); err != nil {
			fmt.Println(err)
			l.Fatal().Err(err).Msg("error running input mode")
		}
		return
	}

	if cfg.Type == output.AppType {
		outputSvc := outputsvc.New(collSvc, metadataSvc, searchSvc)
		if err = output.Run(cfg, outputSvc, l); err != nil {
			fmt.Println(err)
			l.Fatal().Err(err).Msg("error running output mode")
		}
		return
	}
}

package main

import (
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/logger"
	"time"
)

func main() {
	start := time.Now()
	l := logger.New()

	cfg := config.NewConfig()

	if err := cmd.Execute(cfg, l); err != nil {
		l.Fatal().Err(err).Msg("error executing app")
	}

	fmt.Printf("\n completed, took %v\n", time.Since(start))
}

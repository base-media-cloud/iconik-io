package main

import (
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/logger"
	"time"
)

func main() {
	start := time.Now()

	l := logger.New()

	if err := cmd.Execute(l); err != nil {
		l.Fatal().Err(err).Msg("error executing")
	}

	fmt.Printf("\n completed, took %v\n", time.Since(start))
}

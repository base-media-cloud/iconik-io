package main

import (
	"errors"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	logger "github.com/base-media-cloud/pd-iconik-io-rd/internal"
	"time"
)

func main() {
	start := time.Now()
	l, err := logger.New("iconik-io")
	if err != nil {
		l.Fatalw("error encountered: %s", err)
	}
	defer func() {
		err = errors.Join(err, l.Sync())
	}()

	cfg := config.NewConfig()

	if err := cmd.Execute(l, cfg); err != nil {
		l.Fatalw("error encountered:", err)
	}
	fmt.Printf("\n completed, took %v\n", time.Since(start))
}

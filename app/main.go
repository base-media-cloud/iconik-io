package main

import (
	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	logger "github.com/base-media-cloud/pd-iconik-io-rd/internal"
)

func main() {

	l, err := logger.New("iconik-io")
	if err != nil {
		l.Fatalw("error encountered: %s", err)
	}
	defer l.Sync()

	cfg := config.NewConfig()

	if err := cmd.Execute(l, cfg); err != nil {
		l.Fatalw("error encountered:", err)
	}
}

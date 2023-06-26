package main

import (
	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
	"github.com/base-media-cloud/pd-iconik-io-rd/pkg/logger"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

func main() {

	log, err := logger.New("PD-ICONIK-IO-RD")
	if err != nil {
		log.Fatalw("error encountered: %s", err)
	}
	defer log.Sync()

	log.Infow("starting service", zapcore.Field{
		Key:    "build",
		Type:   zapcore.StringType,
		String: build,
	})
	defer log.Infow("shutdown complete")

	if err := cmd.Execute(log); err != nil {
		log.Fatalw("error encountered: %s", err)
	}
}

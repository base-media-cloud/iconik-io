package main

import (
	"log"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error encountered: %s", err)
	}
}

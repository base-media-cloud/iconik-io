/*
Package iconik provides a wrapper for all iconik endpoints.
*/
package iconik

//go:generate mockgen -source api.go -destination=../../mocks/iconik_mocks/api.go -package=iconik_mocks

import (
	"context"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
)

// Requester is an interface that defines the operations that can be performed on an http request.
type Requester interface {
	Do(
		ctx context.Context,
		method,
		url string,
		headers map[string]string,
		queryParams map[string]string,
		payload []byte,
	) ([]byte, *int, error)
}

// API is a wrapper struct for all iconik endpoints.
type API struct {
	cfg     *config.Iconik
	req     Requester
	url     string
	headers map[string]string
}

// New is a function that returns a new instance of the API struct.
func New(cfg *config.Iconik, req Requester) *API {
	return &API{
		cfg: cfg,
		headers: map[string]string{
			"App-ID":     cfg.AppID,
			"Auth-Token": cfg.AuthToken,
		},
		url: cfg.BaseURL,
		req: req,
	}
}

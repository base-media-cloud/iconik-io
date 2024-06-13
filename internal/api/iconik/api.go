/*
Package iconik provides a wrapper for all iconik endpoints.
*/
package iconik

//go:generate mockgen -source api.go -destination=../../mocks/iconik_mocks/api.go -package=iconik_mocks

import (
	"context"

	"golang.org/x/sync/syncmap"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
)

// Requester is an interface that defines the operations that can be performed on an http request.
type Requester interface {
	Do(
		ctx context.Context,
		method,
		url string,
		headers *syncmap.Map,
		queryParams map[string]string,
		payload []byte,
	) ([]byte, *int, error)
}

// API is a wrapper struct for all iconik endpoints.
type API struct {
	cfg     *config.App
	req     Requester
	url     string
	headers *syncmap.Map
}

// New is a function that returns a new instance of the API struct.
func New(cfg *config.App, req Requester) *API {
	headersSyncMap := syncmap.Map{}
	headersSyncMap.Store("App-ID", cfg.AppID)
	headersSyncMap.Store("Auth-Token", cfg.AuthToken)
	return &API{
		cfg:     cfg,
		headers: &headersSyncMap,
		url:     cfg.BaseURL,
		req:     req,
	}
}

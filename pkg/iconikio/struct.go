package iconikio

import (
	"fmt"
	"net/http"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
)

type APIConfig struct {
	Scheme    string
	Host      string
	Endpoints *IconikEndpoints
}

type IconikEndpoints struct {
	Asset        *Endpoints
	Collection   *Endpoints
	MetadataView *Endpoints
	Search       *Endpoints
}

type Endpoints struct {
	Post, Get, Patch, Put, Delete *Endpoint
}

type Endpoint struct {
	Path   string
	Path2  string
	Method string
}

func (c *Client) NewAPIConfig(appCfg config.Config) {
	c.Config.APIConfig = &APIConfig{
		Scheme: "https",
		Host:   c.Config.IconikURL,
		Endpoints: &IconikEndpoints{
			Asset: &Endpoints{
				Get: &Endpoint{
					Path:   appCfg.AssetsPrefixURL,
					Method: http.MethodGet,
				},
				Patch: &Endpoint{
					Path:   appCfg.AssetsPrefixURL,
					Method: http.MethodPatch,
				},
			},
			Collection: &Endpoints{
				Get: &Endpoint{
					Path:   fmt.Sprintf("%s/", appCfg.CollectionPrefixURL),
					Method: http.MethodGet,
				},
			},
			MetadataView: &Endpoints{
				Get: &Endpoint{
					Path:   fmt.Sprintf("%s%s", appCfg.MetadataViewPrefixURL, c.Config.ViewID),
					Method: http.MethodGet,
				},
				Put: &Endpoint{
					Path:   appCfg.SearchPrefixURL,
					Path2:  fmt.Sprintf("/views/%s/", c.Config.ViewID),
					Method: http.MethodPut,
				},
			},
			Search: &Endpoints{
				Post: &Endpoint{
					Path:   appCfg.SearchPrefixURL,
					Method: http.MethodPost,
				},
			},
		},
	}
}

package iconikio

import (
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/config"
	"net/http"
)

type Iconik struct {
	IconikClient *Client
}

type Client struct {
	Assets   []*Asset
	Metadata *Metadata
	Config   *Config
}

// type Conf is the structure that holds the key variables required
// in the execution of the service.
type Config struct {
	IconikURL    string
	AppID        string
	AuthToken    string
	CollectionID string
	ViewID       string
	Input        string
	Output       string
	APIConfig    *APIConfig
}

// Assets is the top level data structure that receives the unmarshalled payload
// response.
type Asset struct {
	Objects []*Object `json:"objects"`
}

// Objects acts as a non nested struct to the Objects type in Assets.
type Object struct {
	ID       string                 `json:"id"`
	Metadata map[string]interface{} `json:"metadata"`
	Title    string                 `json:"title"`
}

// ====================================================
// iconik Objects Response Structure "GET /API/metadata/v1/views/"

// Metadata is the top level data structure that receives the unmarshalled payload
// response.
type Metadata struct {
	ViewFields []*ViewField `json:"view_fields"`
}

// ViewField acts as a non nested struct to the ViewFields type in MetadataFields.
type ViewField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

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

func New(cfg *Config) *Client {
	return &Client{
		Config: cfg,
	}
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
					Path:   fmt.Sprintf("%s%s/contents/", appCfg.CollectionPrefixURL, c.Config.CollectionID),
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

type CSVMetadata struct {
	MetadataValues map[string]struct {
		FieldValues []struct {
			Value string `json:"value"`
		} `json:"field_values"`
	} `json:"metadata_values"`
}

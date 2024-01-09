package iconikio

import (
	"fmt"
	"net/http"

	"github.com/base-media-cloud/pd-iconik-io-rd/config"
)

type Iconik struct {
	IconikClient *Client
}

type Client struct {
	Collection *Collection
	Metadata   *Metadata
	Config     *Config
	Assets     []*Object
}

// Config is the structure that holds the key variables required
// in the execution of the service.
type Config struct {
	IconikURL        string
	AppID            string
	AuthToken        string
	CollectionID     string
	ViewID           string
	Input            string
	Output           string
	APIConfig        *APIConfig
	CSVFilesToUpdate int
	CSVMetadata      []*CSVMetadata
}

// Collection is the top level data structure that receives the unmarshalled payload
// response.
type Collection struct {
	Objects []*Object `json:"objects"`
	Errors  []string  `json:"errors"`
	Pages   int
}

// Object acts as a non nested struct to the Objects type in Collection.
type Object struct {
	ID         string                   `json:"id"`
	Metadata   map[string][]interface{} `json:"metadata"`
	Title      string                   `json:"title"`
	Files      []*File                  `json:"files"`
	ObjectType string                   `json:"object_type"`
}

type File struct {
	DirectoryPath string `json:"directory_path"`
	FileSetId     string `json:"file_set_id"`
	FormatId      string `json:"format_id"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	OriginalName  string `json:"original_name"`
	Size          int    `json:"size"`
	Status        string `json:"status"`
	StorageId     string `json:"storage_id"`
	StorageMethod string `json:"storage_method"`
}

// ====================================================
// iconik Objects Response Structure "GET /API/metadata/v1/views/"

// Metadata is the top level data structure that receives the unmarshalled payload
// response.
type Metadata struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ViewFields  []*ViewField `json:"view_fields"`
	Errors      []string     `json:"errors"`
}

// ViewField acts as a non nested struct to the ViewFields type in Metadata.
type ViewField struct {
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	FieldType string    `json:"field_type"`
	Options   []*Option `json:"options"`
	ReadOnly  bool      `json:"read_only"`
	Required  bool      `json:"required"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
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

type CSVMetadata struct {
	Added                bool
	IDStruct             IDStruct
	OriginalNameStruct   OriginalNameStruct
	SizeStruct           SizeStruct
	TitleStruct          TitleStruct
	MetadataValuesStruct MetadataValuesStruct
}

type IDStruct struct {
	ID string `json:"id"`
}

type OriginalNameStruct struct {
	OriginalName string `json:"original_name"`
}

type SizeStruct struct {
	Size string `json:"size"`
}

type TitleStruct struct {
	Title string `json:"title"`
}

type MetadataValuesStruct struct {
	MetadataValues map[string]struct {
		FieldValues []FieldValue `json:"field_values"`
	} `json:"metadata_values"`
}

type FieldValue struct {
	Value string `json:"value"`
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

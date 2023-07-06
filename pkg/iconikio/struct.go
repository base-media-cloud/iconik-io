package iconikio

type Iconik struct {
	IconikClient *Client
}

type Client struct {
	Assets   []*Asset
	Metadata *Metadata
	Config   *Config
}

// type Conf is the structure that holds the key variaables required
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

type APIConfig struct {
	Scheme    string
	Host      string
	Endpoints map[string]interface{}
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

type IconikApi struct {
	Scheme    string           `yaml:"scheme"`
	Host      string           `yaml:"host"`
	Endpoints *IconikEndpoints `yaml:"endpoints"`
}

type IconikEndpoints struct {
	Asset        []*Endpoint `yaml:"asset"`
	Collection   []*Endpoint `yaml:"collection"`
	MetadataView []*Endpoint `yaml:"metadata-view"`
	Search       []*Endpoint `yaml:"search"`
}

type Endpoint struct {
	Path   []string `yaml:"path"`
	Path2  []string `yaml:"path2"`
	Method string   `yaml:"method"`
}

func New(cfg *Config) *Client {
	return &Client{
		Config: cfg,
	}
}

func (c *Client) NewAPIConfig() {
	c.Config.APIConfig = &APIConfig{
		Scheme: "https",
		Host:   c.Config.IconikURL,
		Endpoints: map[string]interface{}{
			"asset": []interface{}{
				map[string]interface{}{
					"path":   []string{"/API/assets/v1/assets/"},
					"path2":  []string{},
					"method": "GET",
				},
				map[string]interface{}{
					"path":   []string{"/API/assets/v1/assets/"},
					"path2":  []string{},
					"method": "PATCH",
				},
			},
			"collection": []interface{}{
				map[string]interface{}{
					"path":   []string{"/API/assets/v1/collections/", c.Config.CollectionID, "/contents/"},
					"path2":  []string{},
					"method": "GET",
				},
			},
			"metadataView": []interface{}{
				map[string]interface{}{
					"path":   []string{"/API/metadata/v1/views/", c.Config.ViewID},
					"path2":  []string{},
					"method": "GET",
				},
				map[string]interface{}{
					"path":   []string{"/API/metadata/v1/assets/"},
					"path2":  []string{"/views/", c.Config.ViewID, "/"},
					"method": "PUT",
				},
			},
			"search": []interface{}{
				map[string]interface{}{
					"path":   []string{"/API/search/v1/search/"},
					"path2":  []string{},
					"method": "POST",
				},
			},
		},
	}
}

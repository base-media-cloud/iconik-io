package iconikio

type Iconik struct {
	IconikClient *Client
}

type Client struct {
	Assets   []*Asset
	Metadata *Metadata
	cfg      *Config
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

func New(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

// func New(i *Iconik) (*Iconik, error) {

// 	iconikURL := flag.String("iconik-url", "https://preview.iconik.cloud", "iconik URL")
// 	appID := flag.String("app-id", "", "iconik Application ID")
// 	authToken := flag.String("auth-token", "", "iconik Authentication token")
// 	collectionID := flag.String("collection-id", "", "iconik Collection ID")
// 	viewID := flag.String("metadata-view-id", "", "iconik Metadata View ID")
// 	input := flag.String("input", "", "Input mode - requires path to input CSV file")
// 	output := flag.String("output", "", "Output mode - requires path to save CSV file")
// 	flag.Parse()

// 	if *appID == "" {
// 		log.Fatal("No App-Id provided")
// 	}
// 	if *authToken == "" {
// 		log.Fatal("No Auth-Token provided")
// 	}
// 	if *collectionID == "" {
// 		log.Fatal("No Collection ID provided")
// 	}
// 	if *viewID == "" {
// 		log.Fatal("No Metadata View ID provided")
// 	}
// 	if *input == "" && *output == "" {
// 		log.Fatal("Neither input or output mode selected. Please select one.")
// 	}

// 	i.IconikClient.cfg.IconikURL = *iconikURL
// 	i.IconikClient.cfg.AppID = *appID
// 	i.IconikClient.cfg.AuthToken = *authToken
// 	i.IconikClient.cfg.CollectionID = *collectionID
// 	i.IconikClient.cfg.ViewID = *viewID
// 	i.IconikClient.cfg.Input = *input
// 	i.IconikClient.cfg.Output = *output

// 	err := CheckAppIDAuthTokenCollectionID(i.IconikClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = CheckMetadataID(i.IconikClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return i, nil
// }

// func (i *Iconik) Connect() (*Client, error) {

// 	iconikURL := flag.String("iconik-url", "https://preview.iconik.cloud", "iconik URL")
// 	appID := flag.String("app-id", "", "iconik Application ID")
// 	authToken := flag.String("auth-token", "", "iconik Authentication token")
// 	collectionID := flag.String("collection-id", "", "iconik Collection ID")
// 	viewID := flag.String("metadata-view-id", "", "iconik Metadata View ID")
// 	input := flag.String("input", "", "Input mode - requires path to input CSV file")
// 	output := flag.String("output", "", "Output mode - requires path to save CSV file")
// 	flag.Parse()

// 	if *appID == "" {
// 		log.Fatal("No App-Id provided")
// 	}
// 	if *authToken == "" {
// 		log.Fatal("No Auth-Token provided")
// 	}
// 	if *collectionID == "" {
// 		log.Fatal("No Collection ID provided")
// 	}
// 	if *viewID == "" {
// 		log.Fatal("No Metadata View ID provided")
// 	}
// 	if *input == "" && *output == "" {
// 		log.Fatal("Neither input or output mode selected. Please select one.")
// 	}

// 	i.IconikClient.cfg.IconikURL = *iconikURL
// 	i.IconikClient.cfg.AppID = *appID
// 	i.IconikClient.cfg.AuthToken = *authToken
// 	i.IconikClient.cfg.CollectionID = *collectionID
// 	i.IconikClient.cfg.ViewID = *viewID
// 	i.IconikClient.cfg.Input = *input
// 	i.IconikClient.cfg.Output = *output

// 	err := CheckAppIDAuthTokenCollectionID(i.IconikClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = CheckMetadataID(i.IconikClient)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return i.IconikClient, nil
// }

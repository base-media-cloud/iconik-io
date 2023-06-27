package assets

// ====================================================
// iconik Objects Response Structure "GET /v1/assets/"

// Assets is the top level data structure that receives the unmarshalled payload
// response.
type Assets struct {
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

// MetadataFields is the top level data structure that receives the unmarshalled payload
// response.
type MetadataFields struct {
	ViewFields []*ViewField `json:"view_fields"`
}

// ViewField acts as a non nested struct to the ViewFields type in MetadataFields.
type ViewField struct {
	Name string `json:"name"`
}

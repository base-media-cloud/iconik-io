package assets

import "time"

// ==========================================
// iconik Objects Response Structure "GET /v1/assets/"

// Assets is the top level data structure that receives the unmarshalled payload
// response.
type Assets struct {
	FirstURL string    `json:"first_url"`
	LastURL  string    `json:"last_url"`
	NextURL  string    `json:"next_url"`
	Objects  []*Object `json:"objects"`
	Page     int       `json:"page"`
	Pages    int       `json:"pages"`
	PerPage  int       `json:"per_page"`
	PrevURL  string    `json:"prev_url"`
	ScrollID string    `json:"scroll_id"`
	Total    int       `json:"total"`
}

// Objects acts as a non nested struct to the Objects type in Assets.
type Object struct {
	AnalyzeStatus        string    `json:"analyze_status"`
	AncestorCollections  []string  `json:"ancestor_collections"`
	ArchiveStatus        string    `json:"archive_status"`
	Category             string    `json:"category"`
	CreatedByUser        string    `json:"created_by_user"`
	CreatedByUserInfo    *UserInfo `json:"created_by_user_info"`
	CustomKeyframe       string    `json:"custom_keyframe"`
	CustomPoster         string    `json:"custom_poster"`
	DateCreated          string    `json:"date_created"`
	DateDeleted          string    `json:"date_deleted"`
	DateImported         string    `json:"date_imported"`
	DateModified         string    `json:"date_modified"`
	DeletedByUser        string    `json:"deleted_by_user"`
	DeletedByUserInfo    *UserInfo `json:"deleted_by_user_info"`
	DurationMilliseconds int       `json:"duration_milliseconds"`
	ExternalID           string    `json:"external_id"`
	FileNames            []string  `json:"file_names"`
	Files                []*Files  `json:"files"`
	ID                   string    `json:"id"`
	InCollections        []string  `json:"in_collections"`
	IsBlocked            bool      `json:"is_blocked"`
	IsOnline             bool      `json:"is_online"`
	Keyframes            []struct {
	} `json:"keyframes"`
	MediaType string                 `json:"media_type"` // manual addition
	Metadata  map[string]interface{} `json:"metadata"`
	// MetadataView []map[string]interface{} `json:"metadata_view"`
	ObjectType string `json:"object_type"`
	Position   int    `json:"position"`
	Proxies    []struct {
	} `json:"proxies"`
	Relations         []*Relation `json:"relations"`
	Status            string      `json:"status"`
	Title             string      `json:"title"`
	Type              string      `json:"type"`
	UpdatedByUser     string      `json:"updated_by_user"`
	UpdatedByUserInfo *UserInfo   `json:"updated_by_user_info"`
	Versions          []*Version  `json:"versions"`
	VersionsNumber    int         `json:"versions_number"`
	Warning           string      `json:"warning"`
}

// // MetadataView acts as a non nested struct to the MetadataView type in Object.
// type MetadataView struct {
// 	Metadata map[string]interface{}
// }

// Relation acts as a non nested struct to the Relations type in Object.
type Relation struct {
	DateCreated        string `json:"date_created"`
	DateModified       string `json:"date_modified"`
	Description        string `json:"description"`
	RelatedFromAssetID string `json:"related_from_asset_id"`
	RelatedToAssetID   string `json:"related_to_asset_id"`
	RelationType       string `json:"relation_type"`
}

// UserInfo acts as a non nested struct to the recurring UserInfo type in Object.
type UserInfo struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	ID         string `json:"id"`
	LastName   string `json:"last_name"`
	Photo      string `json:"photo"`
	PhotoBig   string `json:"photo_big"`
	PhotoSmall string `json:"photo_small"`
}

// Version acts as a non nested struct to the Version type in Object.
type Version struct {
	AnalyzeStatus     string    `json:"analyze_status"`
	ArchiveStatus     string    `json:"archive_status"`
	CreatedByUser     string    `json:"created_by_user"`
	CreatedByUserInfo *UserInfo `json:"created_by_user_info"`
	DateCreated       string    `json:"date_created"`
	ID                string    `json:"id"`
	IsOnline          bool      `json:"is_online"`
	Status            string    `json:"status"`
	TranscribeStatus  string    `json:"transcribe_status"`
}

// Files acts as a non nested struct to the Files type in Object
type Files struct {
	DirectoryPath string `json:"directory_path"`
	FileSetID     string `json:"file_set_id"`
	FormatID      string `json:"format_id"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	OriginalName  string `json:"original_name"`
	Size          int    `json:"size"`
	Status        string `json:"status"`
	StorageID     string `json:"storage_id"`
	StorageMethod string `json:"storage_method"`
}

type MetadataFields struct {
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
	Description  string    `json:"description"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	ViewFields   []struct {
		AutoSet         bool          `json:"auto_set"`
		DateCreated     time.Time     `json:"date_created"`
		DateModified    time.Time     `json:"date_modified"`
		Description     interface{}   `json:"description"`
		ExternalID      interface{}   `json:"external_id"`
		FieldType       string        `json:"field_type"`
		HideIfNotSet    bool          `json:"hide_if_not_set"`
		IsBlockField    bool          `json:"is_block_field"`
		IsWarningField  bool          `json:"is_warning_field"`
		Label           string        `json:"label"`
		MappedFieldName interface{}   `json:"mapped_field_name"`
		MaxValue        interface{}   `json:"max_value"`
		MinValue        interface{}   `json:"min_value"`
		Multi           bool          `json:"multi"`
		Name            string        `json:"name"`
		Options         []interface{} `json:"options"`
		ReadOnly        bool          `json:"read_only"`
		Representative  bool          `json:"representative"`
		Required        bool          `json:"required"`
		Sortable        bool          `json:"sortable"`
		SourceURL       interface{}   `json:"source_url"`
	} `json:"view_fields"`
}

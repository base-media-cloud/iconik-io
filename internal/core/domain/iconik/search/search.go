package search

import (
	"time"
)

type Search struct {
	DocTypes      []string       `json:"doc_types"`
	Facets        []string       `json:"facets"`
	IncludeFields []string       `json:"include_fields"`
	Sort          []Sort         `json:"sort"`
	Query         string         `json:"query"`
	Filter        Filter         `json:"filter"`
	FacetsFilters []FacetsFilter `json:"facets_filters"`
	SearchFields  []string       `json:"search_fields"`
	SearchAfter   []interface{}  `json:"search_after"`
}

type Sort struct {
	Name  string `json:"name"`
	Order string `json:"order"`
}

type Term struct {
	Name    string   `json:"name"`
	ValueIn []string `json:"value_in"`
}

type FacetsFilter struct {
	Name    string   `json:"name"`
	ValueIn []string `json:"value_in"`
}

type Filter struct {
	Operator string `json:"operator"`
	Terms    []Term `json:"terms"`
}

type Results struct {
	FirstUrl string      `json:"first_url"`
	LastUrl  string      `json:"last_url"`
	NextUrl  string      `json:"next_url"`
	Objects  []Object    `json:"objects"`
	Page     int         `json:"page"`
	Pages    int         `json:"pages"`
	PerPage  int         `json:"per_page"`
	PrevUrl  string      `json:"prev_url"`
	Total    int         `json:"total"`
	Errors   interface{} `json:"errors"`
}

type Object struct {
	Sort                  []interface{}            `json:"_sort"`
	AnalyzeStatus         string                   `json:"analyze_status"`
	AncestorCollections   []string                 `json:"ancestor_collections"`
	ArchiveStatus         string                   `json:"archive_status"`
	Category              interface{}              `json:"category"`
	CreatedByUser         string                   `json:"created_by_user"`
	CreatedByUserInfo     interface{}              `json:"created_by_user_info"`
	DateCreated           time.Time                `json:"date_created"`
	DateModified          time.Time                `json:"date_modified"`
	Duration              string                   `json:"duration"`
	ExternalLink          interface{}              `json:"external_link"`
	Files                 []File                   `json:"files"`
	Format                string                   `json:"format"`
	ID                    string                   `json:"id"`
	InCollections         []string                 `json:"in_collections"`
	IsBlocked             bool                     `json:"is_blocked"`
	IsOnline              bool                     `json:"is_online"`
	MediaType             string                   `json:"media_type"`
	Metadata              map[string][]interface{} `json:"metadata"`
	ObjectType            string                   `json:"object_type"`
	Permissions           []interface{}            `json:"permissions"`
	Position              int                      `json:"position"`
	TimeEndMilliseconds   interface{}              `json:"time_end_milliseconds"`
	TimeStartMilliseconds interface{}              `json:"time_start_milliseconds"`
	Title                 string                   `json:"title"`
	Type                  string                   `json:"type"`
	Versions              []Version                `json:"versions"`
	VersionsNumber        int                      `json:"versions_number"`
	Warning               interface{}              `json:"warning"`
}

// File acts as a non nested struct to the Files type in Object.
type File struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Size         int    `json:"size"`
}

type Version struct {
	AnalyzeStatus    string    `json:"analyze_status"`
	ArchiveStatus    string    `json:"archive_status"`
	CreatedByUser    string    `json:"created_by_user"`
	DateCreated      time.Time `json:"date_created"`
	Id               string    `json:"id"`
	IsOnline         bool      `json:"is_online"`
	Status           string    `json:"status"`
	TranscribeStatus string    `json:"transcribe_status"`
}

// ToResultsDTO is a method that converts a Results to a ResultsDTO.
func (r *Results) ToResultsDTO() ResultsDTO {
	objectDTOs := make([]ObjectDTO, len(r.Objects))
	for i, object := range r.Objects {
		objectDTOs[i] = object.ToObjectDTO()
	}

	return ResultsDTO{
		FirstUrl: r.FirstUrl,
		LastUrl:  r.LastUrl,
		NextUrl:  r.NextUrl,
		Objects:  objectDTOs,
		Page:     r.Page,
		Pages:    r.Pages,
		PerPage:  r.PerPage,
		PrevUrl:  r.PrevUrl,
		Total:    r.Total,
		Errors:   r.Errors,
	}
}

// ToObjectDTO is a method that converts an Object to an ObjectDTO.
func (o *Object) ToObjectDTO() ObjectDTO {
	fileDTOs := make([]FileDTO, len(o.Files))
	for i, file := range o.Files {
		fileDTOs[i] = file.ToFileDTO()
	}

	versionDTOs := make([]VersionDTO, len(o.Versions))
	for i, version := range o.Versions {
		versionDTOs[i] = version.ToVersionDTO()
	}

	return ObjectDTO{
		Sort:                  o.Sort,
		AnalyzeStatus:         o.AnalyzeStatus,
		AncestorCollections:   o.AncestorCollections,
		ArchiveStatus:         o.ArchiveStatus,
		Category:              o.Category,
		CreatedByUser:         o.CreatedByUser,
		CreatedByUserInfo:     o.CreatedByUserInfo,
		DateCreated:           o.DateCreated,
		DateModified:          o.DateModified,
		Duration:              o.Duration,
		ExternalLink:          o.ExternalLink,
		Files:                 fileDTOs,
		Format:                o.Format,
		ID:                    o.ID,
		InCollections:         o.InCollections,
		IsBlocked:             o.IsBlocked,
		IsOnline:              o.IsOnline,
		MediaType:             o.MediaType,
		Metadata:              o.Metadata,
		ObjectType:            o.ObjectType,
		Permissions:           o.Permissions,
		Position:              o.Position,
		TimeEndMilliseconds:   o.TimeEndMilliseconds,
		TimeStartMilliseconds: o.TimeStartMilliseconds,
		Title:                 o.Title,
		Type:                  o.Type,
		Versions:              versionDTOs,
		VersionsNumber:        o.VersionsNumber,
		Warning:               o.Warning,
	}
}

// ToFileDTO is a method that converts a File to a FileDTO.
func (f *File) ToFileDTO() FileDTO {
	return FileDTO{
		Id:           f.Id,
		Name:         f.Name,
		OriginalName: f.OriginalName,
		Size:         f.Size,
	}
}

// ToVersionDTO is a method that converts a Version to a VersionDTO.
func (v *Version) ToVersionDTO() VersionDTO {
	return VersionDTO{
		AnalyzeStatus:    v.AnalyzeStatus,
		ArchiveStatus:    v.ArchiveStatus,
		CreatedByUser:    v.CreatedByUser,
		DateCreated:      v.DateCreated,
		Id:               v.Id,
		IsOnline:         v.IsOnline,
		Status:           v.Status,
		TranscribeStatus: v.TranscribeStatus,
	}
}

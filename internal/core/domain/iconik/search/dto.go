package search

import "time"

type ResultsDTO struct {
	FirstUrl string
	LastUrl  string
	NextUrl  string
	Objects  []ObjectDTO
	Page     int
	Pages    int
	PerPage  int
	PrevUrl  string
	Total    int
	Errors   interface{}
}

type ObjectDTO struct {
	Sort                  []interface{}
	AnalyzeStatus         string
	AncestorCollections   []string
	ArchiveStatus         string
	Category              interface{}
	CreatedByUser         string
	CreatedByUserInfo     interface{}
	DateCreated           time.Time
	DateModified          time.Time
	Duration              string
	ExternalLink          interface{}
	Files                 []FileDTO
	Format                string
	ID                    string
	InCollections         []string
	IsBlocked             bool
	IsOnline              bool
	MediaType             string
	Metadata              map[string][]interface{}
	ObjectType            string
	Permissions           []interface{}
	Position              int
	TimeEndMilliseconds   interface{}
	TimeStartMilliseconds interface{}
	Title                 string
	Type                  string
	Versions              []VersionDTO
	VersionsNumber        int
	Warning               interface{}
}

type FileDTO struct {
	Id           string
	Name         string
	OriginalName string
	Size         int
}

type VersionDTO struct {
	AnalyzeStatus    string
	ArchiveStatus    string
	CreatedByUser    string
	DateCreated      time.Time
	Id               string
	IsOnline         bool
	Status           string
	TranscribeStatus string
}

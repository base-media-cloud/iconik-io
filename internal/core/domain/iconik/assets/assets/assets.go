package assets

import "time"

type Asset struct {
	AnalyzeStatus string    `json:"analyze_status"`
	ArchiveStatus string    `json:"archive_status"`
	CreatedByUser string    `json:"created_by_user"`
	DateCreated   time.Time `json:"date_created"`
	DateImported  time.Time `json:"date_imported"`
	DateModified  time.Time `json:"date_modified"`
	ID            string    `json:"id"`
	IsBlocked     bool      `json:"is_blocked"`
	IsOnline      bool      `json:"is_online"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	UpdatedByUser string    `json:"updated_by_user"`
	Versions      []Version `json:"versions"`
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

// ToDTO is a method that converts a SystemDomain to a DTO.
func (a *Asset) ToDTO() DTO {
	return DTO{
		AnalyzeStatus: a.AnalyzeStatus,
		ArchiveStatus: a.ArchiveStatus,
		CreatedByUser: a.CreatedByUser,
		DateCreated:   a.DateCreated,
		DateImported:  a.DateImported,
		DateModified:  a.DateModified,
		ID:            a.ID,
		IsBlocked:     a.IsBlocked,
		IsOnline:      a.IsOnline,
		Status:        a.Status,
		Title:         a.Title,
		Type:          a.Type,
		UpdatedByUser: a.UpdatedByUser,
	}
}

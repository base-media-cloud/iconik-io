package asset

import "time"

type DTO struct {
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

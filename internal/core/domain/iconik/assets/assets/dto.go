package assets

import "time"

type DTO struct {
	AnalyzeStatus string
	ArchiveStatus string
	CreatedByUser string
	DateCreated   time.Time
	DateImported  time.Time
	DateModified  time.Time
	ID            string
	IsBlocked     bool
	IsOnline      bool
	Status        string
	Title         string
	Type          string
	UpdatedByUser string
	Versions      []Version
}

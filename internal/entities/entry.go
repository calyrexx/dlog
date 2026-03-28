package entities

import "time"

// ValidTags is the set of allowed tag values.
var ValidTags = map[string]bool{
	"feat": true,
	"fix":  true,
	"note": true,
	"idea": true,
	"docs": true,
}

// Entry represents a single diary entry.
type Entry struct {
	ID          int64
	Text        string
	Tag         string
	Repo        string
	Branch      string
	CommitHash  string
	DurationSec int
	CreatedAt   time.Time
}

// StatsResult holds aggregated activity data for a time period.
type StatsResult struct {
	TotalEntries  int
	TotalDuration int // seconds
	ByTag         map[string]int
	ByDay         map[string]int // "2006-01-02" -> count
}

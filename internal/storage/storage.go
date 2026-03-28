package storage

import (
	"time"

	"github.com/calyrexx/dlog/internal/entities"
)

// StatsResult holds aggregated activity data.
type StatsResult struct {
	TotalEntries  int
	TotalDuration int // seconds
	ByTag         map[string]int
	ByDay         map[string]int // "2006-01-02" -> count
}

type Storage interface {
	// Add writes a new entry and returns its ID.
	Add(entry entities.Entry) (int64, error)

	// GetToday returns all entries created today.
	GetToday() ([]entities.Entry, error)

	// GetYesterday returns all entries created yesterday.
	GetYesterday() ([]entities.Entry, error)

	// GetLast returns the n most recent entries.
	GetLast(n int) ([]entities.Entry, error)

	// Search returns entries whose text matches the query.
	Search(query string) ([]entities.Entry, error)

	// GetByDateRange returns entries within [from, to] inclusive.
	GetByDateRange(from, to time.Time) ([]entities.Entry, error)

	// StartSession saves the beginning of a timed work session.
	StartSession(tag, text string) error

	// StopSession finalises the active session, records duration and appends text.
	// Returns the completed entry.
	StopSession(text string) (*entities.Entry, error)

	// ActiveSession returns the in-progress session, or nil if none.
	ActiveSession() (*entities.Entry, error)

	// Stats returns aggregated activity for the given period.
	// period: "day" | "week" | "month" | "year"
	Stats(period string) (*StatsResult, error)

	// Close releases the underlying database connection.
	Close() error
}

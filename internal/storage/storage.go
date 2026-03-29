package storage

import (
	"context"
	"time"

	"github.com/calyrexx/dlog/internal/entities"
)

// Storage defines all data access operations.
type Storage interface {
	// Add writes a new entry and returns its ID.
	Add(ctx context.Context, entry entities.Entry) (int64, error)

	// GetByID returns a single entry by its ID.
	GetByID(ctx context.Context, id int64) (*entities.Entry, error)

	// Update modifies an existing entry's text and tag.
	Update(ctx context.Context, id int64, text, tag string) error

	// Delete removes an entry by its ID.
	Delete(ctx context.Context, id int64) error

	// GetToday returns all entries created today.
	GetToday(ctx context.Context) ([]entities.Entry, error)

	// GetYesterday returns all entries created yesterday.
	GetYesterday(ctx context.Context) ([]entities.Entry, error)

	// GetLast returns the n most recent entries.
	GetLast(ctx context.Context, n int) ([]entities.Entry, error)

	// Search returns entries whose text matches the query.
	Search(ctx context.Context, query, tag, repo string) ([]entities.Entry, error)

	// GetByDateRange returns entries within [from, to] inclusive.
	GetByDateRange(ctx context.Context, from, to time.Time) ([]entities.Entry, error)

	// StartSession saves the beginning of a timed work session.
	StartSession(ctx context.Context, entry entities.Entry) error

	// StopSession finalises the active session, records duration and appends text.
	StopSession(ctx context.Context, text, repo, branch, commitHash string) (*entities.Entry, error)

	// ActiveSession returns the in-progress session.
	// Returns ErrNoActiveSession if no session is running.
	ActiveSession(ctx context.Context) (*entities.Entry, error)

	// Stats returns aggregated activity for the given period.
	Stats(ctx context.Context, period string) (*entities.StatsResult, error)

	// Streaks returns the current and longest consecutive-day streaks.
	Streaks(ctx context.Context) (current, longest int, err error)

	// Close releases the underlying database connection.
	Close() error
}

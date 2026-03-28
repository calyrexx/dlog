package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/calyrexx/dlog/internal/entities"
	_ "modernc.org/sqlite"
)

// ErrNotImplemented is returned by storage methods that are not yet implemented.
var ErrNotImplemented = errors.New("not implemented")

const schema = `
CREATE TABLE IF NOT EXISTS entries (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    text         TEXT NOT NULL,
    tag          TEXT NOT NULL DEFAULT 'note',
    repo         TEXT,
    branch       TEXT,
    commit_hash  TEXT,
    duration_sec INTEGER,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);`

// SQLiteStorage is the SQLite-backed implementation of Storage.
type SQLiteStorage struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

// New opens (or creates) ~/.dlog/dlog.db and applies the schema.
func New() (*SQLiteStorage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	dir := filepath.Join(home, ".dlog")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(dir, "dlog.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if _, err = db.ExecContext(context.Background(), schema); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("init schema: %w", err)
	}

	return &SQLiteStorage{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}, nil
}

func (s *SQLiteStorage) Add(ctx context.Context, e entities.Entry) (int64, error) {
	query, args, err := s.builder.
		Insert("entries").
		Columns("text", "tag", "repo", "branch", "commit_hash", "duration_sec").
		Values(e.Text, e.Tag, e.Repo, e.Branch, e.CommitHash, e.DurationSec).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (s *SQLiteStorage) GetToday(ctx context.Context) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select("id", "text", "tag", "repo", "branch", "commit_hash", "duration_sec", "created_at").
		From("entries").
		Where("date(created_at) = date('now', 'localtime')").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) GetYesterday(ctx context.Context) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select("id", "text", "tag", "repo", "branch", "commit_hash", "duration_sec", "created_at").
		From("entries").
		Where("date(created_at) = date('now', '-1 day', 'localtime')").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) GetLast(ctx context.Context, n int) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select("id", "text", "tag", "repo", "branch", "commit_hash", "duration_sec", "created_at").
		From("entries").
		OrderBy("created_at DESC").
		Limit(uint64(n)). //nolint:gosec
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) Search(ctx context.Context, q string) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select("id", "text", "tag", "repo", "branch", "commit_hash", "duration_sec", "created_at").
		From("entries").
		Where(sq.Like{"text": "%" + q + "%"}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) GetByDateRange(ctx context.Context, from, to time.Time) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select("id", "text", "tag", "repo", "branch", "commit_hash", "duration_sec", "created_at").
		From("entries").
		Where(sq.GtOrEq{"created_at": from}).
		Where(sq.LtOrEq{"created_at": to}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) StartSession(_ context.Context, _, _ string) error {
	return ErrNotImplemented
}

func (s *SQLiteStorage) StopSession(_ context.Context, _ string) (*entities.Entry, error) {
	return nil, ErrNotImplemented
}

func (s *SQLiteStorage) ActiveSession(_ context.Context) (*entities.Entry, error) {
	return nil, ErrNotImplemented
}

func (s *SQLiteStorage) Stats(_ context.Context, _ string) (*StatsResult, error) {
	return nil, ErrNotImplemented
}

func (s *SQLiteStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	return nil
}

// queryEntries executes a SELECT and scans the result into []Entry.
func (s *SQLiteStorage) queryEntries(ctx context.Context, query string, args ...any) ([]entities.Entry, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}
	defer rows.Close()

	var entries []entities.Entry

	for rows.Next() {
		var e entities.Entry
		if err := rows.Scan(
			&e.ID, &e.Text, &e.Tag, &e.Repo, &e.Branch, &e.CommitHash, &e.DurationSec, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return entries, nil
}

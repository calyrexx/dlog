package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/calyrexx/dlog/internal/entities"
	_ "modernc.org/sqlite"
)

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
	dir := filepath.Join(os.Getenv("HOME"), ".dlog")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(dir, "dlog.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if _, err = db.Exec(schema); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("init schema: %w", err)
	}

	return &SQLiteStorage{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}, nil
}

func (s *SQLiteStorage) Add(e entities.Entry) (int64, error) {
	query, args, err := s.builder.
		Insert("entries").
		Columns(
			"text",
			"tag",
			"repo",
			"branch",
			"commit_hash",
			"duration_sec",
		).
		Values(
			e.Text,
			e.Tag,
			e.Repo,
			e.Branch,
			e.CommitHash,
			e.DurationSec,
		).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	res, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (s *SQLiteStorage) GetToday() ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select(
			"id",
			"text",
			"tag",
			"repo",
			"branch",
			"commit_hash",
			"duration_sec",
			"created_at",
		).
		From("entries").
		Where("date(created_at) = date('now', 'localtime')").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}
	defer rows.Close()

	var entries []entities.Entry

	for rows.Next() {
		var e entities.Entry
		if err := rows.Scan(
			&e.ID,
			&e.Text,
			&e.Tag,
			&e.Repo,
			&e.Branch,
			&e.CommitHash,
			&e.DurationSec,
			&e.CreatedAt,
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

func (s *SQLiteStorage) GetYesterday() ([]entities.Entry, error) {
	// TODO: SELECT WHERE date(created_at) = date('now', '-1 day', 'localtime')
	return nil, nil
}

func (s *SQLiteStorage) GetLast(n int) ([]entities.Entry, error) {
	// TODO: SELECT ORDER BY created_at DESC LIMIT n
	return nil, nil
}

func (s *SQLiteStorage) Search(_ string) ([]entities.Entry, error) {
	// TODO: SELECT WHERE text LIKE '%query%' ORDER BY created_at DESC
	return nil, nil
}

func (s *SQLiteStorage) GetByDateRange(_, _ time.Time) ([]entities.Entry, error) {
	// TODO: SELECT WHERE created_at BETWEEN from AND to
	return nil, nil
}

func (s *SQLiteStorage) StartSession(_, _ string) error {
	// TODO: persist session start to a sidecar file ~/.devlog/session.json
	// or a dedicated in-memory/DB row with duration_sec = -1 as sentinel
	return nil
}

func (s *SQLiteStorage) StopSession(_ string) (*entities.Entry, error) {
	// TODO: read active session, compute duration_sec, INSERT entry, remove session file
	return nil, nil
}

func (s *SQLiteStorage) ActiveSession() (*entities.Entry, error) {
	// TODO: read session file / sentinel row; return nil if none exists
	return nil, nil
}

func (s *SQLiteStorage) Stats(_ string) (*StatsResult, error) {
	// TODO: aggregate entries for period using GROUP BY date/tag
	return nil, nil
}

func (s *SQLiteStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	return nil
}

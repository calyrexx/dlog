package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	db *sql.DB
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

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Add(_ entities.Entry) (int64, error) {
	// TODO: INSERT INTO entries and return lastInsertId
	return 0, nil
}

func (s *SQLiteStorage) GetToday() ([]entities.Entry, error) {
	// TODO: SELECT WHERE date(created_at) = date('now', 'localtime')
	return nil, nil
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
	return s.db.Close()
}

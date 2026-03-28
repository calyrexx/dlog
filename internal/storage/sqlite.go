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

var (
	// ErrSessionActive is returned when starting a session while one is already active.
	ErrSessionActive = errors.New("a session is already active")
	// ErrNoActiveSession is returned when no session is running.
	ErrNoActiveSession = errors.New("no active session")
	// ErrNotFound is returned when the requested entry does not exist.
	ErrNotFound = errors.New("entry not found")
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

var entryCols = []string{
	"id", "text", "tag",
	"COALESCE(repo, '')", "COALESCE(branch, '')",
	"COALESCE(commit_hash, '')",
	"COALESCE(duration_sec, 0)", "created_at",
}

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

func (s *SQLiteStorage) GetByID(ctx context.Context, id int64) (*entities.Entry, error) {
	query, args, err := s.builder.
		Select(entryCols...).
		From("entries").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var e entities.Entry

	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&e.ID, &e.Text, &e.Tag, &e.Repo, &e.Branch,
		&e.CommitHash, &e.DurationSec, &e.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("scan row: %w", err)
	}

	e.CreatedAt = toLocal(e.CreatedAt)

	return &e, nil
}

func (s *SQLiteStorage) Update(ctx context.Context, id int64, text, tag string) error {
	query, args, err := s.builder.
		Update("entries").
		Set("text", text).
		Set("tag", tag).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SQLiteStorage) Delete(ctx context.Context, id int64) error {
	query, args, err := s.builder.
		Delete("entries").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *SQLiteStorage) GetToday(ctx context.Context) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select(entryCols...).
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
		Select(entryCols...).
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
	if n <= 0 {
		return nil, nil
	}

	query, args, err := s.builder.
		Select(entryCols...).
		From("entries").
		OrderBy("created_at DESC").
		Limit(uint64(n)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) Search(ctx context.Context, q, tag, repo string) ([]entities.Entry, error) {
	b := s.builder.
		Select(entryCols...).
		From("entries")

	if q != "" {
		b = b.Where(sq.Like{"text": "%" + q + "%"})
	}

	if tag != "" {
		b = b.Where(sq.Eq{"tag": tag})
	}

	if repo != "" {
		b = b.Where(sq.Eq{"repo": repo})
	}

	query, args, err := b.
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return s.queryEntries(ctx, query, args...)
}

func (s *SQLiteStorage) GetByDateRange(ctx context.Context, from, to time.Time) ([]entities.Entry, error) {
	query, args, err := s.builder.
		Select(entryCols...).
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

func (s *SQLiteStorage) StartSession(ctx context.Context, e entities.Entry) error {
	_, err := s.ActiveSession(ctx)
	if err == nil {
		return ErrSessionActive
	}

	if !errors.Is(err, ErrNoActiveSession) {
		return fmt.Errorf("check active session: %w", err)
	}

	query, args, buildErr := s.builder.
		Insert("entries").
		Columns("text", "tag", "repo", "branch", "commit_hash", "duration_sec").
		Values(e.Text, e.Tag, e.Repo, e.Branch, e.CommitHash, -1).
		ToSql()
	if buildErr != nil {
		return fmt.Errorf("build query: %w", buildErr)
	}

	if _, err = s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s *SQLiteStorage) StopSession(
	ctx context.Context, text, repo, branch, commitHash string,
) (*entities.Entry, error) {
	active, err := s.ActiveSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("check active session: %w", err)
	}

	elapsed := int(time.Since(active.CreatedAt).Seconds())
	if elapsed < 0 {
		elapsed = 0
	}

	finalText := active.Text
	if text != "" {
		if finalText != "" {
			finalText += " — " + text
		} else {
			finalText = text
		}
	}

	ub := s.builder.
		Update("entries").
		Set("duration_sec", elapsed).
		Set("text", finalText)

	if repo != "" {
		ub = ub.Set("repo", repo).
			Set("branch", branch).
			Set("commit_hash", commitHash)
	}

	query, args, err := ub.Where(sq.Eq{"id": active.ID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	if _, err = s.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	active.DurationSec = elapsed
	active.Text = finalText

	if repo != "" {
		active.Repo = repo
		active.Branch = branch
		active.CommitHash = commitHash
	}

	return active, nil
}

func (s *SQLiteStorage) ActiveSession(ctx context.Context) (*entities.Entry, error) {
	query, args, err := s.builder.
		Select(entryCols...).
		From("entries").
		Where(sq.Eq{"duration_sec": -1}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var e entities.Entry

	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&e.ID, &e.Text, &e.Tag, &e.Repo, &e.Branch,
		&e.CommitHash, &e.DurationSec, &e.CreatedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNoActiveSession
	case err != nil:
		return nil, fmt.Errorf("scan row: %w", err)
	}

	e.CreatedAt = toLocal(e.CreatedAt)

	return &e, nil
}

func (s *SQLiteStorage) Stats(ctx context.Context, period string) (*entities.StatsResult, error) {
	from, to := periodRange(period)

	query, args, err := s.builder.
		Select(entryCols...).
		From("entries").
		Where(sq.GtOrEq{"created_at": from}).
		Where(sq.LtOrEq{"created_at": to}).
		Where(sq.NotEq{"duration_sec": -1}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	entries, err := s.queryEntries(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result := &entities.StatsResult{
		TotalEntries: len(entries),
		ByTag:        make(map[string]int),
		ByDay:        make(map[string]int),
	}

	for _, e := range entries {
		result.TotalDuration += e.DurationSec
		result.ByTag[e.Tag]++
		result.ByDay[e.CreatedAt.Format("2006-01-02")]++
	}

	return result, nil
}

func periodRange(period string) (time.Time, time.Time) {
	now := time.Now()
	to := now

	var from time.Time

	switch period {
	case "day":
		y, m, d := now.Date()
		from = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	case "week":
		offset := int(now.Weekday()+6) % 7
		y, m, d := now.AddDate(0, 0, -offset).Date()
		from = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	case "month":
		y, m, _ := now.Date()
		from = time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
	case "year":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		offset := int(now.Weekday()+6) % 7
		y, m, d := now.AddDate(0, 0, -offset).Date()
		from = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	}

	return from, to
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
			&e.ID, &e.Text, &e.Tag, &e.Repo, &e.Branch,
			&e.CommitHash, &e.DurationSec, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		e.CreatedAt = toLocal(e.CreatedAt)

		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return entries, nil
}

// toLocal re-interprets a UTC-parsed time as local time.
// SQLite stores datetime('now','localtime') without timezone info,
// so the Go driver parses it as UTC.
func toLocal(t time.Time) time.Time {
	if t.Location() == time.UTC {
		return time.Date(
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
			t.Nanosecond(), time.Now().Location(),
		)
	}

	return t
}

-- +goose Up
CREATE TABLE IF NOT EXISTS entries (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    text         TEXT NOT NULL,
    tag          TEXT NOT NULL DEFAULT 'note',
    repo         TEXT,
    issue_link   TEXT,
    branch       TEXT,
    commit_hash  TEXT,
    duration_sec INTEGER,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);
-- +goose Down
DROP TABLE IF EXISTS entries;

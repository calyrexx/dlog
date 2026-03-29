# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build     # Build binary to ./dlog  (injects version via ldflags)
make install   # Install binary to $GOPATH/bin
make test      # Run tests: go test ./...
make lint      # Run golangci-lint
make checks    # Run test + lint together
make stage m="commit message"  # Run checks, commit, and push current branch
```

Run a single test:
```bash
go test ./internal/storage/... -run TestFunctionName
```

## Architecture

**dlog** is a developer diary CLI (Cobra) that logs timestamped entries to SQLite (`~/.dlog/dlog.db`), automatically capturing git context (repo, branch, commit hash).

```
main.go → command.App
  command/*.go       CLI handlers (add, today, yesterday, last, week, search, start/stop, status, edit, delete, export, stats)
  storage/           Storage interface + SQLite implementation (Squirrel query builder)
  entities/entry.go  Entry + StatsResult data models, ValidTags set
  git/git.go         Extract repo name, branch, commit from shell commands
  render/render.go   Terminal output: tables (lipgloss), contribution graph, bar charts, status messages
```

**Dependency injection:** `command.App` holds a `storage.Storage` interface, injected at startup in `main.go`. Commands receive the app struct — not the storage directly.

**Entry flow:** `add cmd` → `git.*` (extract context) → `storage.Add()` → SQLite

**Session tracking:** active sessions are entries with `duration_sec = -1`. `StartSession` inserts such a row; `StopSession` calculates elapsed time and updates it to the actual duration. Only one session can be active at a time.

**Timezone handling:** SQLite stores `datetime('now','localtime')` without timezone info. The Go driver parses it as UTC, so `toLocal()` in `sqlite.go` re-interprets scanned times in the local timezone. This must be applied after every row scan (queryEntries, GetByID, ActiveSession).

## Key Implementation Notes

- Linter is strict: 120-char line limit, max cyclomatic complexity 15, max func length 200 lines / 50 statements, `paralleltest` enforced. Zero `nolint` directives — restructure code instead.
- Valid tags are defined in `entities.ValidTags`: `feat`, `fix`, `note`, `idea`, `docs`. Tag validation runs in commands via `validateTag()` in `root.go`. Tag colors live in `render/render.go`.
- Git functions return empty strings on failure (non-git directories), not errors. Commands work fine outside git repos.
- `render` depends only on `entities`, never on `storage` — keep this direction clean.
- Default command (`dlog` with no args) shows today's entries.
- `ContributionGraph` adapts its width based on the stats period (1/8/16/52 weeks).
- Valid stats periods: `day`, `week` (default), `month`, `year`. `PeriodRange` / `PrevPeriodRange` in `sqlite.go` compute the time boundaries used by both `Stats` and the comparison panel.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build     # Build binary to ./dlog
make install   # Install binary to $GOPATH/bin
make test      # Run tests: go test ./...
make lint      # Run golangci-lint
make checks    # Run test + lint together
```

Run a single test:
```bash
go test ./internal/storage/... -run TestFunctionName
```

## Architecture

**dlog** is a developer diary CLI (Cobra) that logs timestamped entries to SQLite (`~/.dlog/dlog.db`), automatically capturing git context (repo, branch, commit hash).

```
main.go → command.App
  command/*.go       CLI handlers (add, today; others are stubs)
  storage/           Storage interface + SQLite implementation (Squirrel query builder)
  entities/entry.go  Entry data model
  git/git.go         Extract repo name, branch, commit from shell commands
  render/render.go   Terminal output: tables (lipgloss), contribution graph
```

**Dependency injection:** `command.App` holds a `storage.Storage` interface, injected at startup in `main.go`. Commands receive the app struct — not the storage directly.

**Entry flow:** `add cmd` → `git.*` (extract context) → `storage.Add()` → SQLite

## Key Implementation Notes

- Linter is strict: 120-char line limit, max cyclomatic complexity 15, max func length 200 lines / 50 statements, `paralleltest` enforced.
- Storage methods `StartSession`, `StopSession`, `ActiveSession`, and `Stats` are defined on the interface but not yet implemented in `sqlite.go`.
- Commands `yesterday`, `last`, `search`, `start`, `stop`, `export`, `stats` are registered but return stubs — implement storage methods first, then the command handler, then the render output.
- Tag color mapping lives in `render/render.go`; valid tags are: `feat`, `fix`, `note`, `idea`, `docs`.
- Git functions execute shell commands via `exec.CommandContext` — they return empty strings on failure (non-git directories), not errors.

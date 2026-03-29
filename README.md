# dlog

Developer diary for the terminal. Log what you did, when, and why - with automatic git context.
## Install

```bash
go install github.com/calyrexx/dlog@latest
```

Or build from source:

```bash
make build    # ./dlog
make install  # $GOPATH/bin/dlog
```

## Quick start

```bash
dlog add "implemented auth middleware" -t feat
dlog add "fixed nil pointer in handler" -t fix
dlog add "need to revisit caching logic" -t idea
dlog # show today's entries
```

## Commands

### Writing entries

```bash
dlog add <text> [-t tag]       # add an entry (default tag: note)
dlog edit <id> [text] [-t tag] # edit text or tag
dlog delete <id>               # remove an entry
```

Tags: `feat` `fix` `note` `idea` `docs`

### Timed sessions

```bash
dlog start [text] [-t tag]    # start a work session
dlog status                   # check if a session is running
dlog stop [text]              # stop and save with duration
```

### Viewing entries

```bash
dlog                          # today (default)
dlog today                    # same as above
dlog yesterday                # yesterday's entries
dlog week                     # current week (Mon-Sun)
dlog last [-n 20]             # N most recent (default: 10)
```

### Search

```bash
dlog search <query>           # full-text search
dlog search -t feat           # all entries with tag
dlog search auth -t fix       # combine query + tag
```

### Stats & export

```bash
dlog stats [-p day|week|month|year]  # activity stats + contribution graph
dlog export [-f json|csv|markdown] [-o file] [--from 2025-01-01] [--to 2025-12-31]
```

## Git context

When run inside a git repository, `dlog add` and `dlog start/stop` automatically capture the repo name, branch, and commit hash. Outside git - entries are saved without git context.

## Data

All entries are stored in `~/.dlog/dlog.db` (SQLite).

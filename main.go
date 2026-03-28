package main

import (
	"log/slog"
	"os"

	"github.com/calyrexx/dlog/internal/command"
	"github.com/calyrexx/dlog/internal/storage"
	"github.com/calyrexx/zeroslog"
)

func main() {
	logger := slog.New(zeroslog.New(
		zeroslog.WithOutput(os.Stderr),
		zeroslog.WithMinLevel(slog.LevelDebug),
		zeroslog.WithColors(),
	))
	slog.SetDefault(logger)

	db, err := storage.New()
	if err != nil {
		slog.Error("failed to open storage", "err", err)
		os.Exit(1)
	}
	defer func(db *storage.SQLiteStorage) {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close storage", "err", err)
		}
	}(db)

	app := command.NewApp(db)
	app.Execute()
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := storage.New()
	if err != nil {
		return fmt.Errorf("open storage: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close storage", "err", err)
		}
	}()

	app := command.NewApp(db)
	app.Execute(ctx)

	return nil
}

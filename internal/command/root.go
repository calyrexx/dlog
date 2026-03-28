package command

import (
	"context"
	"os"

	"github.com/calyrexx/dlog/internal/storage"
	"github.com/spf13/cobra"
)

type App struct {
	db storage.Storage
}

func NewApp(db storage.Storage) *App {
	return &App{db: db}
}

// Execute builds the command tree and runs it.
func (a *App) Execute(ctx context.Context) {
	root := &cobra.Command{
		Use:   "dlog",
		Short: "Personal developer diary for the terminal",
		Long:  "dlog — track what you did, when you did it, and why.",
	}

	root.AddCommand(
		a.newAddCmd(),
		a.newTodayCmd(),
		a.newLastCmd(),
		a.newSearchCmd(),
		a.newStartCmd(),
		a.newStopCmd(),
		a.newYesterdayCmd(),
		a.newExportCmd(),
		a.newStatsCmd(),
	)

	root.SetContext(ctx)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

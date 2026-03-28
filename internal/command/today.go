package command

import (
	"fmt"
	"log/slog"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newTodayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "today",
		Short: "Show all entries created today",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("today command", "args", args)

			notes, err := a.db.GetToday()
			if err != nil {
				slog.Error("today command", "error", err)

				return fmt.Errorf("get today db error: %w", err)
			}

			slog.Debug("today command", slog.Group("notes", "count", len(notes)))

			render.Table(notes)

			return nil
		},
	}
}

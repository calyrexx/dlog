package command

import (
	"fmt"
	"log/slog"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newLastCmd() *cobra.Command {
	var n int

	cmd := &cobra.Command{
		Use:   "last",
		Short: "Show the N most recent entries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			slog.Debug("last command", "count", n)

			entries, err := a.db.GetLast(cmd.Context(), n)
			if err != nil {
				slog.Error("last command", "error", err)

				return fmt.Errorf("get last db error: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}

	cmd.Flags().IntVarP(&n, "count", "n", 10, "number of entries to show")

	return cmd
}

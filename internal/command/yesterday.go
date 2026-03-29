package command

import (
	"fmt"
	"log/slog"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newYesterdayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "yesterday",
		Short: "Show all entries created yesterday",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("yesterday command")

			entries, err := a.db.GetYesterday(cmd.Context())
			if err != nil {
				slog.Error("yesterday command", "error", err)

				return fmt.Errorf("get yesterday db error: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}
}

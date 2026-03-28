package command

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newStatsCmd() *cobra.Command {
	var period string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show activity statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			slog.Debug("stats command", "period", period)

			result, err := a.db.Stats(ctx, period)
			if err != nil {
				slog.Error("stats command", "error", err)

				return fmt.Errorf("stats db error: %w", err)
			}

			render.Stats(result, period)

			// contribution graph: fetch last 26 weeks
			now := time.Now()
			graphFrom := now.AddDate(0, 0, -52*7)

			entries, err := a.db.GetByDateRange(ctx, graphFrom, now)
			if err != nil {
				slog.Error("stats command", "graph error", err)

				return fmt.Errorf("graph db error: %w", err)
			}

			counts := make(map[string]int, len(entries))
			for _, e := range entries {
				counts[e.CreatedAt.Format("2006-01-02")]++
			}

			render.ContributionGraph(counts)

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&period, "period", "p", "week",
		"time period: day | week | month | year",
	)

	return cmd
}

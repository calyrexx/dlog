package command

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

var graphRangeDays = map[string]int{
	"day":   7,
	"week":  7,
	"month": 4 * 7,
	"year":  52 * 7,
}

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
				return fmt.Errorf("stats: %w", err)
			}

			render.Stats(result, period)

			now := time.Now()

			days := graphRangeDays[period]
			if days == 0 {
				days = graphRangeDays["week"]
			}

			entries, err := a.db.GetByDateRange(ctx, now.AddDate(0, 0, -days), now)
			if err != nil {
				return fmt.Errorf("graph: %w", err)
			}

			counts := make(map[string]int, len(entries))
			for _, e := range entries {
				counts[e.CreatedAt.Format("2006-01-02")]++
			}

			render.ContributionGraph(counts, period)

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&period, "period", "p", "week",
		"time period: day | week | month | year",
	)

	return cmd
}

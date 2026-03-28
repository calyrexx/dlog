package command

import (
	"fmt"
	"log/slog"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/calyrexx/dlog/internal/storage"
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
				return fmt.Errorf("stats: %w", err)
			}

			cur, best, streakErr := a.db.Streaks(ctx)
			if streakErr != nil {
				return fmt.Errorf("streaks: %w", streakErr)
			}

			prevFrom, prevTo := storage.PrevPeriodRange(period)

			prevEntries, err := a.db.GetByDateRange(ctx, prevFrom, prevTo)
			if err != nil {
				return fmt.Errorf("prev period: %w", err)
			}

			var prevDuration int

			for _, e := range prevEntries {
				if e.DurationSec > 0 {
					prevDuration += e.DurationSec
				}
			}

			prev := &render.PrevPeriod{
				Entries:  len(prevEntries),
				Duration: prevDuration,
			}

			from, to := storage.PeriodRange(period)
			render.Stats(result, period, cur, best, prev, from, to)

			entries, err := a.db.GetByDateRange(ctx, from, to)
			if err != nil {
				return fmt.Errorf("graph: %w", err)
			}

			counts := make(map[string]int, len(entries))
			hourly := make(map[int]map[int]int)

			for _, e := range entries {
				counts[e.CreatedAt.Format("2006-01-02")]++

				wd := (int(e.CreatedAt.Weekday()) + 6) % 7 // Monday=0
				h := e.CreatedAt.Hour()

				if hourly[wd] == nil {
					hourly[wd] = make(map[int]int)
				}

				hourly[wd][h]++
			}

			render.ContributionGraph(counts, period)
			render.HourlyHeatmap(hourly)

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&period, "period", "p", "week",
		"time period: day | week | month | year",
	)

	return cmd
}

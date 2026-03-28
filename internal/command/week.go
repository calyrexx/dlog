package command

import (
	"fmt"
	"time"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newWeekCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "week",
		Short: "Show all entries from this week (Mon-Sun)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			now := time.Now()
			offset := int(now.Weekday()+6) % 7
			y, m, d := now.AddDate(0, 0, -offset).Date()

			weekStart := time.Date(y, m, d, 0, 0, 0, 0, now.Location())

			entries, err := a.db.GetByDateRange(cmd.Context(), weekStart, now)
			if err != nil {
				return fmt.Errorf("get week entries: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}
}

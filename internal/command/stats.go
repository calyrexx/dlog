package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newStatsCmd() *cobra.Command {
	var period string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show activity statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")

			return nil
		},
	}

	cmd.Flags().StringVarP(&period, "period", "p", "week", "time period: day | week | month | year")

	return cmd
}

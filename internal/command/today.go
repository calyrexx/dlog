package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newTodayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "today",
		Short: "Show all entries created today",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")
			return nil
		},
	}
}

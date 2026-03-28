package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newYesterdayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "yesterday",
		Short: "Show all entries created yesterday",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")
			return nil
		},
	}
}

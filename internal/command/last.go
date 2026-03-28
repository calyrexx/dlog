package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newLastCmd() *cobra.Command {
	var n int

	cmd := &cobra.Command{
		Use:   "last",
		Short: "Show the N most recent entries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")
			return nil
		},
	}

	cmd.Flags().IntVarP(&n, "count", "n", 10, "number of entries to show")

	return cmd
}

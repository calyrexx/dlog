package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newSearchCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search across all entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "filter results by tag")

	return cmd
}

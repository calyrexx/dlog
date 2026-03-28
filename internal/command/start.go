package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newStartCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "start [text]",
		Short: "Start a timed work session",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note", "tag for the session entry")

	return cmd
}

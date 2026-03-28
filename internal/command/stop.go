package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [text]",
		Short: "Stop the active timed session and save the entry",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")

			return nil
		},
	}
}

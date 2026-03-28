package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *App) newExportCmd() *cobra.Command {
	var (
		format string
		output string
		from   string
		to     string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export entries to a file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("not implemented yet")

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "output format: json | csv | markdown")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path (default: stdout)")
	cmd.Flags().StringVar(&from, "from", "", "start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&to, "to", "", "end date (YYYY-MM-DD)")

	return cmd
}

package command

import (
	"fmt"
	"strconv"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an entry by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			if err := a.db.Delete(cmd.Context(), id); err != nil {
				return fmt.Errorf("delete entry: %w", err)
			}

			render.EntryDeleted(id)

			return nil
		},
	}
}

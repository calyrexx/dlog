package command

import (
	"fmt"
	"log/slog"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newSearchCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search across all entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("search command", "query", args[0], "tag", tag)

			entries, err := a.db.Search(cmd.Context(), args[0])
			if err != nil {
				slog.Error("search command", "error", err)

				return fmt.Errorf("search db error: %w", err)
			}

			if tag != "" {
				entries = filterByTag(entries, tag)
			}

			render.Table(entries)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "filter results by tag")

	return cmd
}

func filterByTag(entries []entities.Entry, tag string) []entities.Entry {
	filtered := make([]entities.Entry, 0, len(entries))

	for _, e := range entries {
		if e.Tag == tag {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

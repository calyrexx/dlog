package command

import (
	"fmt"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newSearchCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search entries by text and/or tag",
		Long: "Search entries by text query, tag filter, or both.\n\n" +
			"Examples:\n  dlog search auth\n  dlog search -t feat\n" +
			"  dlog search migration -t fix",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			query := ""

			if len(args) > 0 {
				query = args[0]
			}

			if query == "" && tag == "" {
				return fmt.Errorf("provide a query, a --tag filter, or both")
			}

			var (
				entries []entities.Entry
				err     error
			)

			switch {
			case query != "" && tag != "":
				entries, err = a.db.Search(ctx, query)
				if err == nil {
					entries = filterByTag(entries, tag)
				}
			case tag != "":
				entries, err = a.db.SearchByTag(ctx, tag)
			default:
				entries, err = a.db.Search(ctx, query)
			}

			if err != nil {
				return fmt.Errorf("search: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "filter by tag")

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

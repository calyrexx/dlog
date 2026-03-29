package command

import (
	"fmt"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newSearchCmd() *cobra.Command {
	var (
		tag  string
		repo string
	)

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search entries by text and/or tag",
		Long: "Search entries by text query, project tag filter, or both.\n\n" +
			"Examples:\n  dlog search auth\n  dlog search -t feat\n" +
			"  dlog search migration -t fix",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var query string

			if len(args) > 0 {
				query = args[0]
			}

			if query == "" && tag == "" && repo == "" {
				return fmt.Errorf("provide a query, a --tag filter, a --project filter, or all")
			}

			entries, err := a.db.Search(ctx, query, tag, repo)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "filter by tag")
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "filter by git repository")

	return cmd
}

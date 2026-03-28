package command

import (
	"fmt"
	"strings"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/git"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newAddCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "add <text>",
		Short: "Add a new diary entry",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTag(tag); err != nil {
				return err
			}

			ctx := cmd.Context()
			text := strings.Join(args, " ")

			id, err := a.db.Add(ctx,
				entities.Entry{
					Text:       text,
					Tag:        tag,
					Repo:       git.RepoName(ctx),
					Branch:     git.Branch(ctx),
					CommitHash: git.CommitHash(ctx),
				},
			)
			if err != nil {
				return fmt.Errorf("add entry: %w", err)
			}

			render.EntryAdded(id, tag, text)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note",
		fmt.Sprintf("entry tag (%s)", ValidTagList))

	return cmd
}

package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newEditCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "edit <id> [new text]",
		Short: "Edit an entry's text or tag",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}

			ctx := cmd.Context()

			entry, err := a.db.GetByID(ctx, id)
			if err != nil {
				return fmt.Errorf("get entry: %w", err)
			}

			newText := entry.Text
			if len(args) > 1 {
				newText = strings.Join(args[1:], " ")
			}

			newTag := entry.Tag

			if cmd.Flags().Changed("tag") {
				if err := validateTag(tag); err != nil {
					return err
				}

				newTag = tag
			}

			if err := a.db.Update(ctx, id, newText, newTag); err != nil {
				return fmt.Errorf("update entry: %w", err)
			}

			entry.Text = newText
			entry.Tag = newTag

			render.EntryUpdated(entry)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "",
		fmt.Sprintf("new tag (%s)", ValidTagList))

	return cmd
}

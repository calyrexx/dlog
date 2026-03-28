package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/git"
	"github.com/spf13/cobra"
)

func (a *App) newAddCmd() *cobra.Command {
	var (
		tag      string
		noteText string
	)

	cmd := &cobra.Command{
		Use:   "add <text>",
		Short: "Add a new diary entry",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("add command", "args", args)

			noteText = strings.Join(args, " ")

			repo, err := git.RepoName()
			if err != nil {
				slog.Error("add command", "error", err)

				return fmt.Errorf("failed to get: %w", err)
			}

			branch, err := git.Branch()
			if err != nil {
				slog.Error("add command", "error", err)

				return fmt.Errorf("failed to get: %w", err)
			}

			commitHash, err := git.CommitHash()
			if err != nil {
				slog.Error("add command", "error", err)

				return fmt.Errorf("failed to get: %w", err)
			}

			id, err := a.db.Add(entities.Entry{
				Text:       noteText,
				Tag:        tag,
				Repo:       repo,
				Branch:     branch,
				CommitHash: commitHash,
			})
			if err != nil {
				slog.Error("add command", "error", err)

				return fmt.Errorf("add entry db error: %w", err)
			}

			slog.Debug("add command", slog.Group("note", "id", id, "text", noteText, "tag", tag))

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note", "entry tag (note, fix, feat, idea, …)")

	return cmd
}

package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/git"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/calyrexx/dlog/internal/storage"
	"github.com/spf13/cobra"
)

func (a *App) newStartCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "start [text]",
		Short: "Start a timed work session",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTag(tag); err != nil {
				return err
			}

			ctx := cmd.Context()
			text := strings.Join(args, " ")

			err := a.db.StartSession(ctx, entities.Entry{
				Text:       text,
				Tag:        tag,
				Repo:       git.RepoName(ctx),
				Branch:     git.Branch(ctx),
				CommitHash: git.CommitHash(ctx),
			})

			if errors.Is(err, storage.ErrSessionActive) {
				active, aErr := a.db.ActiveSession(ctx)
				if aErr != nil {
					return fmt.Errorf("start session: %w", err)
				}

				render.ActiveSessionStatus(active)
				fmt.Println()
				fmt.Println("  stop it first: dlog stop")

				return nil
			}

			if err != nil {
				return fmt.Errorf("start session: %w", err)
			}

			render.SessionStarted(tag, text)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note",
		fmt.Sprintf("tag for the session (%s)", ValidTagList))

	return cmd
}

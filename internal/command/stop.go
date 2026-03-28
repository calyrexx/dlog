package command

import (
	"fmt"
	"strings"

	"github.com/calyrexx/dlog/internal/git"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [text]",
		Short: "Stop the active timed session and save the entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			text := strings.Join(args, " ")

			entry, err := a.db.StopSession(ctx, text,
				git.RepoName(ctx), git.Branch(ctx), git.CommitHash(ctx))
			if err != nil {
				return fmt.Errorf("stop session: %w", err)
			}

			render.SessionStopped(entry)

			return nil
		},
	}
}

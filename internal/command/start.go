package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newStartCmd() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "start [text]",
		Short: "Start a timed work session",
		RunE: func(cmd *cobra.Command, args []string) error {
			text := strings.Join(args, " ")

			slog.Debug("start command", "tag", tag, "text", text)

			if err := a.db.StartSession(cmd.Context(), tag, text); err != nil {
				slog.Error("start command", "error", err)

				return fmt.Errorf("start session: %w", err)
			}

			render.SessionStarted(tag, text)

			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note", "tag for the session entry")

	return cmd
}

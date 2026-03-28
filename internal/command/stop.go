package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/spf13/cobra"
)

func (a *App) newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [text]",
		Short: "Stop the active timed session and save the entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			text := strings.Join(args, " ")

			slog.Debug("stop command", "text", text)

			entry, err := a.db.StopSession(cmd.Context(), text)
			if err != nil {
				slog.Error("stop command", "error", err)

				return fmt.Errorf("stop session: %w", err)
			}

			render.SessionStopped(entry)

			return nil
		},
	}
}

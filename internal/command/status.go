package command

import (
	"errors"
	"fmt"

	"github.com/calyrexx/dlog/internal/render"
	"github.com/calyrexx/dlog/internal/storage"
	"github.com/spf13/cobra"
)

func (a *App) newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the active session, if any",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			session, err := a.db.ActiveSession(cmd.Context())
			if errors.Is(err, storage.ErrNoActiveSession) {
				render.NoActiveSession()

				return nil
			}

			if err != nil {
				return fmt.Errorf("active session: %w", err)
			}

			render.ActiveSessionStatus(session)

			return nil
		},
	}
}

package command

import (
	"log/slog"
	"strings"

	"github.com/calyrexx/dlog/internal/entities"
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
			return nil
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "note", "entry tag (note, fix, feat, idea, …)")

	id, err := a.db.Add(entities.Entry{
		Text:        noteText,
		Tag:         tag,
		Repo:        "",
		Branch:      "",
		CommitHash:  "",
		DurationSec: 0,
	})
	if err != nil {
		slog.Error("add command", "error", err)
	}

	slog.Debug("add command", "note.id", id, "note.text", noteText)

	return cmd
}

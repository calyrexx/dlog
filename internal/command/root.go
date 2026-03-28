package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/render"
	"github.com/calyrexx/dlog/internal/storage"
	"github.com/spf13/cobra"
)

// ValidTagList returns a comma-separated list of valid tags for help text.
var ValidTagList = func() string {
	tags := make([]string, 0, len(entities.ValidTags))
	for t := range entities.ValidTags {
		tags = append(tags, t)
	}

	return strings.Join(tags, ", ")
}()

type App struct {
	db storage.Storage
}

func NewApp(db storage.Storage) *App {
	return &App{db: db}
}

// Execute builds the command tree and runs it.
func (a *App) Execute(ctx context.Context) {
	root := &cobra.Command{
		Use:   "dlog",
		Short: "Personal developer diary for the terminal",
		Long:  "dlog — track what you did, when you did it, and why.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := a.db.GetToday(cmd.Context())
			if err != nil {
				return fmt.Errorf("get today: %w", err)
			}

			render.Table(entries)

			return nil
		},
	}

	root.AddCommand(
		a.newAddCmd(),
		a.newTodayCmd(),
		a.newYesterdayCmd(),
		a.newWeekCmd(),
		a.newLastCmd(),
		a.newSearchCmd(),
		a.newStartCmd(),
		a.newStopCmd(),
		a.newStatusCmd(),
		a.newDeleteCmd(),
		a.newEditCmd(),
		a.newExportCmd(),
		a.newStatsCmd(),
	)

	root.SetContext(ctx)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func validateTag(tag string) error {
	if !entities.ValidTags[tag] {
		return fmt.Errorf("unknown tag %q, valid tags: %s", tag, ValidTagList)
	}

	return nil
}

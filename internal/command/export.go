package command

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/spf13/cobra"
)

func (a *App) newExportCmd() *cobra.Command {
	var (
		format string
		output string
		from   string
		to     string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export entries to a file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			slog.Debug("export command",
				"format", format, "output", output,
				"from", from, "to", to,
			)

			fromTime, toTime, err := parseDateRange(from, to)
			if err != nil {
				return err
			}

			entries, err := a.db.GetByDateRange(ctx, fromTime, toTime)
			if err != nil {
				return fmt.Errorf("get entries: %w", err)
			}

			var w io.Writer = os.Stdout

			if output != "" {
				f, fErr := os.Create(output)
				if fErr != nil {
					return fmt.Errorf("create file: %w", fErr)
				}

				defer f.Close()

				w = f
			}

			switch format {
			case "json":
				return exportJSON(w, entries)
			case "csv":
				return exportCSV(w, entries)
			case "markdown", "md":
				return exportMarkdown(w, entries)
			default:
				return fmt.Errorf("unknown format: %s", format)
			}
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json",
		"output format: json | csv | markdown")
	cmd.Flags().StringVarP(&output, "output", "o", "",
		"output file path (default: stdout)")
	cmd.Flags().StringVar(&from, "from", "",
		"start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&to, "to", "",
		"end date (YYYY-MM-DD)")

	return cmd
}

func parseDateRange(from, to string) (time.Time, time.Time, error) {
	var (
		fromTime time.Time
		toTime   time.Time
		err      error
	)

	if from != "" {
		fromTime, err = time.ParseInLocation("2006-01-02", from,
			time.Now().Location())
		if err != nil {
			return fromTime, toTime,
				fmt.Errorf("invalid --from date: %w", err)
		}
	} else {
		fromTime = time.Date(2000, 1, 1, 0, 0, 0, 0,
			time.Now().Location())
	}

	if to != "" {
		toTime, err = time.ParseInLocation("2006-01-02", to,
			time.Now().Location())
		if err != nil {
			return fromTime, toTime,
				fmt.Errorf("invalid --to date: %w", err)
		}

		toTime = toTime.Add(24*time.Hour - time.Second)
	} else {
		toTime = time.Now()
	}

	return fromTime, toTime, nil
}

type jsonEntry struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Tag       string `json:"tag"`
	Repo      string `json:"repo,omitempty"`
	Branch    string `json:"branch,omitempty"`
	Commit    string `json:"commit,omitempty"`
	Duration  int    `json:"duration_sec,omitempty"`
	CreatedAt string `json:"created_at"`
}

func exportJSON(w io.Writer, entries []entities.Entry) error {
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = jsonEntry{
			ID:        e.ID,
			Text:      e.Text,
			Tag:       e.Tag,
			Repo:      e.Repo,
			Branch:    e.Branch,
			Commit:    e.CommitHash,
			Duration:  e.DurationSec,
			CreatedAt: e.CreatedAt.Format(time.RFC3339),
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func exportCSV(w io.Writer, entries []entities.Entry) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{
		"id", "text", "tag", "repo", "branch",
		"commit", "duration_sec", "created_at",
	}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}

	for _, e := range entries {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.Text,
			e.Tag,
			e.Repo,
			e.Branch,
			e.CommitHash,
			fmt.Sprintf("%d", e.DurationSec),
			e.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}

	return nil
}

func exportMarkdown(w io.Writer, entries []entities.Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No entries found.")

		return err //nolint:wrapcheck
	}

	currentDate := ""

	for _, e := range entries {
		date := e.CreatedAt.Format("2006-01-02")
		if date != currentDate {
			if currentDate != "" {
				fmt.Fprintln(w)
			}

			fmt.Fprintf(w, "## %s\n\n", date)

			currentDate = date
		}

		timeStr := e.CreatedAt.Format(time.TimeOnly)
		fmt.Fprintf(w, "- **[%s]** `%s` %s", timeStr, e.Tag, e.Text)

		if e.Repo != "" {
			fmt.Fprintf(w, " (%s/%s)", e.Repo, e.Branch)
		}

		fmt.Fprintln(w)
	}

	return nil
}

package render

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/calyrexx/dlog/internal/entities"
	"github.com/calyrexx/dlog/internal/storage"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	colorSubtle = lipgloss.Color("240")
	colorAccent = lipgloss.Color("99")
	colorMuted  = lipgloss.Color("243")

	styleHeader = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			Padding(0, 1)

	styleBorder = lipgloss.NewStyle().
			Foreground(colorSubtle)

	styleCell = lipgloss.NewStyle().
			Padding(0, 1)

	styleTag = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)
)

var tagColors = map[string]lipgloss.Color{
	"feat": "10",
	"fix":  "9",
	"note": "12",
	"idea": "11",
	"docs": "14",
}

var (
	graphLevels = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("237")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("22")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("28")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("34")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("46")),
	}

	graphDayLabels = [7]string{"Mon", "", "Wed", "", "Fri", "", ""}
)

func Table(entries []entities.Entry) {
	if len(entries) == 0 {
		fmt.Println("no entries found")

		return
	}

	rows := make([][]string, len(entries))
	for i, e := range entries {
		rows[i] = []string{
			e.CreatedAt.Format(time.TimeOnly),
			e.Tag,
			e.Text,
			e.Repo,
			e.Branch,
		}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(styleBorder).
		Headers("TIME", "TAG", "TEXT", "REPO", "BRANCH").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styleHeader
			}

			if col == 1 { // TAG column
				tag := rows[row][1]
				if c, ok := tagColors[tag]; ok {
					return lipgloss.NewStyle().
						Foreground(c).
						Padding(0, 1)
				}

				return styleTag
			}

			return styleCell
		})

	fmt.Println(t.Render())
}

const (
	graphWeeks     = 52
	graphCellWidth = 2
	graphGap       = 1
	graphColWidth  = graphCellWidth + graphGap
	graphLabelW    = 6
)

var graphCell = strings.Repeat("█", graphCellWidth)

// ContributionGraph renders a GitHub-style activity heatmap from a day->count map.
func ContributionGraph(counts map[string]int) {
	today := time.Now()
	todayWd := int(today.Weekday()+6) % 7
	start := today.AddDate(0, 0, -(todayWd + (graphWeeks-1)*7 + 5))

	total := 0
	for _, v := range counts {
		total += v
	}

	muted := lipgloss.NewStyle().Foreground(colorMuted)
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))

	fmt.Println()
	fmt.Printf("  %s in the last year\n\n",
		title.Render(fmt.Sprintf("%d contributions", total)))

	graphMonthLabels(start, muted)
	graphRows(start, today, counts, muted)
	graphLegend(muted)
}

func graphMonthLabels(start time.Time, muted lipgloss.Style) {
	labels := make([]string, graphWeeks)
	lastMonth := -1

	for w := range graphWeeks {
		m := int(start.AddDate(0, 0, w*7).Month())
		if m != lastMonth {
			labels[w] = start.AddDate(0, 0, w*7).Format("Jan")
			lastMonth = m
		}
	}

	var line strings.Builder

	line.WriteString(strings.Repeat(" ", graphLabelW))

	for w := range graphWeeks {
		if labels[w] != "" {
			// skip label if next month label is too close (< 2 cols apart)
			hasSpace := w+1 >= graphWeeks || labels[w+1] == "" //nolint:gosec // guarded
			if hasSpace {
				line.WriteString(labels[w])

				target := graphLabelW + (w+1)*graphColWidth
				for line.Len() < target {
					line.WriteByte(' ')
				}

				continue
			}
		}

		target := graphLabelW + (w+1)*graphColWidth
		for line.Len() < target {
			line.WriteByte(' ')
		}
	}

	fmt.Println(muted.Render(line.String()))
}

func graphRows(start, today time.Time, counts map[string]int, muted lipgloss.Style) {
	for row := range 7 {
		label := graphDayLabels[row]
		if label == "" {
			label = "   "
		}

		fmt.Print(muted.Render(fmt.Sprintf("  %s ", label)))

		for w := range graphWeeks {
			day := start.AddDate(0, 0, w*7+row)
			if day.After(today) {
				fmt.Print(strings.Repeat(" ", graphColWidth))

				continue
			}

			n := counts[day.Format("2006-01-02")]
			fmt.Print(graphLevels[level(n)].Render(graphCell))
			fmt.Print(strings.Repeat(" ", graphGap))
		}

		fmt.Println()
	}
}

func graphLegend(muted lipgloss.Style) {
	fmt.Println()

	gridWidth := graphLabelW + graphWeeks*graphColWidth

	var buf strings.Builder

	buf.WriteString(muted.Render("Less "))

	for _, s := range graphLevels {
		buf.WriteString(s.Render(graphCell))
		buf.WriteString(" ")
	}

	buf.WriteString(muted.Render("More"))

	visibleLen := 5 + 5*graphCellWidth + 4 + 4
	pad := gridWidth - visibleLen

	if pad < 0 {
		pad = 0
	}

	fmt.Printf("%s%s\n", strings.Repeat(" ", pad), buf.String())
}

// BarChart prints a horizontal bar chart of string->count data to stdout.
func BarChart(data map[string]int) {
	if len(data) == 0 {
		return
	}

	type kv struct {
		key   string
		value int
	}

	sorted := make([]kv, 0, len(data))
	maxVal := 0
	maxKeyLen := 0

	for k, v := range data {
		sorted = append(sorted, kv{k, v})
		if v > maxVal {
			maxVal = v
		}

		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})

	const maxBarWidth = 30

	barStyle := lipgloss.NewStyle().Foreground(colorAccent)
	labelStyle := lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(maxKeyLen).
		Align(lipgloss.Right)

	for _, item := range sorted {
		width := maxBarWidth
		if maxVal > 0 {
			width = item.value * maxBarWidth / maxVal
		}

		if width == 0 && item.value > 0 {
			width = 1
		}

		bar := barStyle.Render(strings.Repeat("█", width))
		label := labelStyle.Render(item.key)
		fmt.Printf("  %s %s %d\n", label, bar, item.value)
	}
}

// Stats prints formatted statistics output.
func Stats(result *storage.StatsResult, period string) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorAccent).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255")).
		MarginTop(1)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	fmt.Println(titleStyle.Render(
		fmt.Sprintf("  Activity — %s", period),
	))

	fmt.Printf("  Entries:  %s\n", valueStyle.Render(
		fmt.Sprintf("%d", result.TotalEntries),
	))

	if result.TotalDuration > 0 {
		fmt.Printf("  Duration: %s\n",
			valueStyle.Render(formatDuration(result.TotalDuration)))
	}

	fmt.Printf("  Days:     %s\n", valueStyle.Render(
		fmt.Sprintf("%d", len(result.ByDay)),
	))

	if len(result.ByTag) > 0 {
		fmt.Println(headerStyle.Render("  By tag"))
		BarChart(result.ByTag)
	}
}

// SessionStarted prints a confirmation when a session starts.
func SessionStarted(tag, text string) {
	accent := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	muted := lipgloss.NewStyle().Foreground(colorMuted)

	fmt.Printf(
		"  %s session started %s\n",
		accent.Render("▶"),
		muted.Render(time.Now().Format(time.TimeOnly)),
	)

	if text != "" {
		fmt.Printf("  %s %s\n",
			muted.Render("text:"), text)
	}

	fmt.Printf("  %s %s\n", muted.Render("tag: "), tagStyled(tag))
}

// SessionStopped prints a confirmation when a session stops.
func SessionStopped(e *entities.Entry) {
	accent := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	muted := lipgloss.NewStyle().Foreground(colorMuted)

	fmt.Printf(
		"  %s session stopped — %s\n",
		accent.Render("■"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("10")).
			Render(formatDuration(e.DurationSec)),
	)

	if e.Text != "" {
		fmt.Printf("  %s %s\n", muted.Render("text:"), e.Text)
	}
}

func tagStyled(tag string) string {
	if c, ok := tagColors[tag]; ok {
		return lipgloss.NewStyle().Foreground(c).Render(tag)
	}

	return lipgloss.NewStyle().Foreground(colorMuted).Render(tag)
}

// EntryAdded prints a confirmation when an entry is added.
func EntryAdded(id int64, tag, text string) {
	accent := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	muted := lipgloss.NewStyle().Foreground(colorMuted)

	fmt.Printf("  %s entry #%d added  %s  %s\n",
		accent.Render("✓"),
		id,
		tagStyled(tag),
		muted.Render(text),
	)
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60

	switch {
	case h > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case m > 0:
		return fmt.Sprintf("%dm", m)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}

// level maps an entry count to a palette index.
func level(n int) int {
	switch {
	case n == 0:
		return 0
	case n <= 2:
		return 1
	case n <= 5:
		return 2
	case n <= 9:
		return 3
	default:
		return 4
	}
}

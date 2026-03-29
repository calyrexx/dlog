package render

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/calyrexx/dlog/internal/entities"
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

const tagCol = 2

// Table renders entries as a styled table with ID, time, tag, text, and git info.
func Table(entries []entities.Entry) {
	if len(entries) == 0 {
		fmt.Println("  no entries found")

		return
	}

	rows := make([][]string, len(entries))
	for i, e := range entries {
		rows[i] = []string{
			fmt.Sprintf("#%d", e.ID),
			e.CreatedAt.Format(time.TimeOnly),
			e.Tag,
			e.Text,
			e.Repo,
			e.Branch,
		}
	}

	styleID := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Padding(0, 1)

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(styleBorder).
		Headers("ID", "TIME", "TAG", "TEXT", "REPO", "BRANCH").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styleHeader
			}

			if col == 0 {
				return styleID
			}

			if col == tagCol {
				tag := rows[row][tagCol]
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

var (
	accent = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	muted  = lipgloss.NewStyle().Foreground(colorMuted)
)

// ActiveSessionStatus prints info about the currently running session.
func ActiveSessionStatus(e *entities.Entry) {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	elapsed := int(time.Since(e.CreatedAt).Seconds())

	fmt.Printf("  %s active session\n", accent.Render("▶"))
	fmt.Printf("  %s %s\n",
		muted.Render("elapsed:"),
		green.Render(formatDuration(elapsed)))
	fmt.Printf("  %s %s\n",
		muted.Render("started:"),
		e.CreatedAt.Format("15:04:05 (2006-01-02)"))
	fmt.Printf("  %s %s\n", muted.Render("tag:    "), tagStyled(e.Tag))

	if e.Text != "" {
		fmt.Printf("  %s %s\n", muted.Render("text:   "), e.Text)
	}

	if e.Repo != "" {
		fmt.Printf("  %s %s/%s\n", muted.Render("git:    "), e.Repo, e.Branch)
	}
}

// NoActiveSession prints a message when no session is running.
func NoActiveSession() {
	fmt.Printf("  %s\n", muted.Render("no active session"))
}

// EntryDeleted prints a confirmation when an entry is deleted.
func EntryDeleted(id int64) {
	fmt.Printf("  %s entry #%d deleted\n", accent.Render("✓"), id)
}

// EntryUpdated prints a confirmation when an entry is updated.
func EntryUpdated(e *entities.Entry) {
	fmt.Printf("  %s entry #%d updated  %s  %s\n",
		accent.Render("✓"),
		e.ID,
		tagStyled(e.Tag),
		muted.Render(e.Text),
	)
}

const (
	graphCellWidth = 2
	graphGap       = 1
	graphColWidth  = graphCellWidth + graphGap
	graphLabelW    = 6
)

var (
	graphCell = strings.Repeat("█", graphCellWidth)

	periodWeeks = map[string]int{
		"day":   1,
		"week":  2,
		"month": 5,
		"year":  52,
	}

	periodTitles = map[string]string{
		"day":   "today",
		"week":  "the last week",
		"month": "the last month",
		"year":  "the last year",
	}
)

// ContributionGraph renders a GitHub-style activity heatmap from a day->count map.
func ContributionGraph(counts map[string]int, period string) {
	weeks := periodWeeks[period]
	if weeks == 0 {
		weeks = periodWeeks["week"]
	}

	today := time.Now()
	todayWd := int(today.Weekday()+6) % 7
	start := today.AddDate(0, 0, -(todayWd + (weeks-1)*7 + 5))

	total := 0
	for _, v := range counts {
		total += v
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))

	label := periodTitles[period]
	if label == "" {
		label = periodTitles["week"]
	}

	fmt.Println()
	fmt.Printf("  %s in %s\n\n",
		title.Render(fmt.Sprintf("%d notes", total)), label)

	graphMonthLabels(start, weeks, muted)
	graphRows(start, today, weeks, counts, muted)
	graphLegend(weeks, muted)
}

func graphMonthLabels(start time.Time, weeks int, muted lipgloss.Style) {
	labels := make([]string, weeks)
	lastMonth := -1

	for w := range weeks {
		m := int(start.AddDate(0, 0, w*7).Month())
		if m != lastMonth {
			labels[w] = start.AddDate(0, 0, w*7).Format("Jan")
			lastMonth = m
		}
	}

	var line strings.Builder

	line.WriteString(strings.Repeat(" ", graphLabelW))

	for w := range weeks {
		if labels[w] != "" && hasLabelSpace(labels, w) {
			line.WriteString(labels[w])

			target := graphLabelW + (w+1)*graphColWidth
			for line.Len() < target {
				line.WriteByte(' ')
			}

			continue
		}

		target := graphLabelW + (w+1)*graphColWidth
		for line.Len() < target {
			line.WriteByte(' ')
		}
	}

	fmt.Println(muted.Render(line.String()))
}

func graphRows(start, today time.Time, weeks int, counts map[string]int, muted lipgloss.Style) {
	for row := range 7 {
		label := graphDayLabels[row]
		if label == "" {
			label = "   "
		}

		fmt.Print(muted.Render(fmt.Sprintf("  %s ", label)))

		for w := range weeks {
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

func graphLegend(weeks int, muted lipgloss.Style) {
	fmt.Println()

	gridWidth := graphLabelW + weeks*graphColWidth

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

// PrevPeriod holds aggregated data for the previous period comparison.
type PrevPeriod struct {
	Entries  int
	Duration int // seconds
}

type StatsArgs struct {
	Result        *entities.StatsResult
	Period        string
	CurrentStreak int
	LongestStreak int
	Prev          *PrevPeriod
	From          time.Time
	To            time.Time
}

// Stats prints formatted statistics output.
func Stats(args StatsArgs) {
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

	labelStyle := lipgloss.NewStyle().Foreground(colorMuted).Width(12).Align(lipgloss.Right)

	rangeStr := formatRange(args.From, args.To)
	fmt.Println(titleStyle.Render(
		fmt.Sprintf("  Activity — %s  %s", args.Period, muted.Render(rangeStr)),
	))

	fmt.Printf("  %s  %s%s\n",
		labelStyle.Render("Entries"),
		valueStyle.Render(fmt.Sprintf("%d", args.Result.TotalEntries)),
		deltaStr(args.Result.TotalEntries, args.Prev.Entries, ""),
	)

	if args.Result.TotalDuration > 0 {
		fmt.Printf("  %s  %s%s\n",
			labelStyle.Render("Duration"),
			valueStyle.Render(formatDuration(args.Result.TotalDuration)),
			deltaStr(args.Result.TotalDuration, args.Prev.Duration, "duration"),
		)
	}

	activeDays := len(args.Result.ByDay)
	fmt.Printf("  %s  %s\n",
		labelStyle.Render("Active days"),
		valueStyle.Render(fmt.Sprintf("%d", activeDays)),
	)

	if activeDays > 0 {
		avg := float64(args.Result.TotalEntries) / float64(activeDays)
		fmt.Printf("  %s  %s\n",
			labelStyle.Render("Avg/day"),
			valueStyle.Render(fmt.Sprintf("%.1f", avg)),
		)
	}

	if args.LongestStreak > 0 {
		streakText := fmt.Sprintf("%d days", args.CurrentStreak)
		if args.LongestStreak > args.CurrentStreak {
			streakText += fmt.Sprintf(" (best: %d)", args.LongestStreak)
		}

		fmt.Printf("  %s  %s\n",
			labelStyle.Render("Streak"),
			valueStyle.Render(streakText),
		)
	}

	if best, count := mostActiveDay(args.Result.ByDay); best != "" {
		t, parseErr := time.Parse("2006-01-02", best)

		label := best
		if parseErr == nil {
			label = t.Format("Mon, Jan 2")
		}

		fmt.Printf("  %s  %s %s\n",
			labelStyle.Render("Peak day"),
			valueStyle.Render(label),
			fmt.Sprintf("(%d entries)", count),
		)
	}

	if len(args.Result.ByTag) > 0 {
		fmt.Println()
		fmt.Println(headerStyle.Render("  By tag"))
		tagBarChart(args.Result.ByTag, args.Result.TotalEntries)
		fmt.Println()
	}

	if len(args.Result.ByRepo) > 0 {
		fmt.Println(headerStyle.Render("  By project"))
		projectBranches(args.Result.ByRepo, args.Result.RepoBranches)
	}
}

func formatRange(from, to time.Time) string {
	if from.Year() == to.Year() && from.Month() == to.Month() && from.Day() == to.Day() {
		return from.Format("Jan 2, 2006")
	}

	if from.Year() == to.Year() {
		return fmt.Sprintf("%s – %s", from.Format("Jan 2"), to.Format("Jan 2"))
	}

	return fmt.Sprintf("%s – %s", from.Format("Jan 2, 2006"), to.Format("Jan 2, 2006"))
}

func mostActiveDay(byDay map[string]int) (string, int) {
	best := ""
	peak := 0

	for day, count := range byDay {
		if count > peak {
			peak = count
			best = day
		}
	}

	return best, peak
}

func tagBarChart(data map[string]int, total int) {
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

	const maxBarWidth = 25

	labelStyle := lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(maxKeyLen).
		Align(lipgloss.Right)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(4).
		Align(lipgloss.Right)

	for _, item := range sorted {
		width := maxBarWidth
		if maxVal > 0 {
			width = item.value * maxBarWidth / maxVal
		}

		if width == 0 && item.value > 0 {
			width = 1
		}

		barColor := colorAccent
		if c, ok := tagColors[item.key]; ok {
			barColor = c
		}

		barStyle := lipgloss.NewStyle().Foreground(barColor)
		bar := barStyle.Render(strings.Repeat("█", width))
		label := labelStyle.Render(item.key)

		pct := 0
		if total > 0 {
			pct = item.value * 100 / total
		}

		fmt.Printf("  %s %s %s %s\n",
			label, bar,
			countStyle.Render(fmt.Sprintf("%d", item.value)),
			fmt.Sprintf("%d%%", pct),
		)
	}
}

func projectBranches(repos map[string]int, branches map[string]map[string]int) {
	type kv struct {
		key   string
		value int
	}

	sorted := make([]kv, 0, len(repos))

	maxKeyLen := 0
	total := 0

	for k, v := range repos {
		sorted = append(sorted, kv{k, v})
		total += v

		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})

	limit := len(sorted)
	if limit > 5 {
		limit = 5
	}

	maxKeyLen = branchLabelWidth(branches, maxKeyLen)

	const maxBarWidth = 25

	maxVal := sorted[0].value
	barStyle := lipgloss.NewStyle().Foreground(colorAccent)
	repoLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true).
		Width(maxKeyLen).Align(lipgloss.Right)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(4).
		Align(lipgloss.Right)

	for i, repo := range sorted[:limit] {
		if i > 0 {
			fmt.Println()
		}

		width := repo.value * maxBarWidth / maxVal
		if width == 0 && repo.value > 0 {
			width = 1
		}

		pct := 0
		if total > 0 {
			pct = repo.value * 100 / total
		}

		bar := barStyle.Render(strings.Repeat("█", width))
		fmt.Printf("  %s %s %s %s\n",
			repoLabelStyle.Render(repo.key), bar,
			countStyle.Render(fmt.Sprintf("%d", repo.value)),
			fmt.Sprintf("%d%%", pct),
		)

		if br, ok := branches[repo.key]; ok && len(br) > 1 {
			renderBranches(br, maxVal, maxBarWidth, maxKeyLen)
		}
	}
}

func branchLabelWidth(branches map[string]map[string]int, w int) int {
	for _, br := range branches {
		for k := range br {
			if len(k)+2 > w {
				w = len(k) + 2
			}
		}
	}

	return w
}

func renderBranches(br map[string]int, maxVal, maxBarWidth, labelWidth int) {
	type kv struct {
		key   string
		value int
	}

	brSorted := make([]kv, 0, len(br))
	for k, v := range br {
		brSorted = append(brSorted, kv{k, v})
	}

	sort.Slice(brSorted, func(i, j int) bool {
		return brSorted[i].value > brSorted[j].value
	})

	brLimit := len(brSorted)
	if brLimit > 3 {
		brLimit = 3
	}

	branchBarStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	branchLabelStyle := lipgloss.NewStyle().Foreground(colorMuted).
		Width(labelWidth).Align(lipgloss.Right)

	countStyle := lipgloss.NewStyle().
		Foreground(colorMuted).
		Width(4).
		Align(lipgloss.Right)

	for _, b := range brSorted[:brLimit] {
		bw := b.value * maxBarWidth / maxVal
		if bw == 0 && b.value > 0 {
			bw = 1
		}

		bBar := branchBarStyle.Render(strings.Repeat("░", bw))
		fmt.Printf("  %s %s %s\n",
			branchLabelStyle.Render(b.key), bBar,
			countStyle.Render(fmt.Sprintf("%d", b.value)),
		)
	}
}

// SessionStarted prints a confirmation when a session starts.
func SessionStarted(tag, text string) {
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

func hasLabelSpace(labels []string, w int) bool {
	next := w + 1
	if next >= len(labels) {
		return true
	}

	return labels[next] == ""
}

// HourlyHeatmap renders a 7x24 grid showing activity by weekday and hour.
func HourlyHeatmap(data map[int]map[int]int) {
	total := 0

	for _, hours := range data {
		for _, n := range hours {
			total += n
		}
	}

	if total == 0 {
		return
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))

	fmt.Println()
	fmt.Printf("  %s\n\n", title.Render("Activity by hour"))

	heatmapHourLabels()

	days := [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for wd := range 7 {
		fmt.Print(muted.Render(fmt.Sprintf("  %s ", days[wd])))

		for h := range 24 {
			n := 0
			if data[wd] != nil {
				n = data[wd][h]
			}

			fmt.Print(graphLevels[level(n)].Render(graphCell))
			fmt.Print(" ")
		}

		fmt.Println()
	}

	fmt.Println()
}

func heatmapHourLabels() {
	var line strings.Builder

	line.WriteString("      ")

	for h := range 24 {
		if h%3 == 0 {
			label := fmt.Sprintf("%02d", h)
			line.WriteString(label)

			target := 6 + (h+1)*graphColWidth
			for line.Len() < target {
				line.WriteByte(' ')
			}
		} else {
			target := 6 + (h+1)*graphColWidth
			for line.Len() < target {
				line.WriteByte(' ')
			}
		}
	}

	fmt.Println(muted.Render(line.String()))
}

func deltaStr(current, previous int, mode string) string {
	if previous == 0 {
		return ""
	}

	diff := current - previous

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	style := green
	arrow := " +"

	if diff < 0 {
		style = red
		arrow = " -"
	} else if diff == 0 {
		return muted.Render("  =")
	}

	absDiff := diff
	if absDiff < 0 {
		absDiff = -absDiff
	}

	var text string
	if mode == "duration" {
		text = fmt.Sprintf("%s%s", arrow, formatDuration(absDiff))
	} else {
		text = fmt.Sprintf("%s%d", arrow, absDiff)
	}

	return "  " + style.Render(text)
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

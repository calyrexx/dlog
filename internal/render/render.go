package render

import (
	"fmt"
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
		lipgloss.NewStyle().Background(lipgloss.Color("237")), // пусто
		lipgloss.NewStyle().Background(lipgloss.Color("22")),  // тёмно-зелёный
		lipgloss.NewStyle().Background(lipgloss.Color("28")),  // зелёный
		lipgloss.NewStyle().Background(lipgloss.Color("34")),  // яркий
		lipgloss.NewStyle().Background(lipgloss.Color("46")),  // максимум
	}

	dayLabels = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
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

func ContributionGraph(entries []entities.Entry) {
	const (
		weeks = 26
		cell  = "  "
	)

	// bucket: "2006-01-02" -> count
	counts := make(map[string]int, len(entries))
	for _, e := range entries {
		counts[e.CreatedAt.Format("2006-01-02")]++
	}

	// начало сетки: первый понедельник >= (сегодня - 26 недель)
	today := time.Now()
	// откатываемся до ближайшего понедельника включительно
	offset := int(today.Weekday()+6) % 7 // 0=Mon … 6=Sun
	startMonday := today.AddDate(0, 0, -(offset + (weeks-1)*7))

	// month labels над графиком
	var monthLine strings.Builder
	monthLine.WriteString("    ") // отступ под day labels

	lastMonth := -1

	for w := range weeks {
		day := startMonday.AddDate(0, 0, w*7)
		if int(day.Month()) != lastMonth {
			label := day.Format("Jan")
			monthLine.WriteString(label)

			lastMonth = int(day.Month())
			// каждая ячейка = 2 символа; label занял 3 — добираем пробел
			monthLine.WriteString(" ")
		} else {
			monthLine.WriteString("   ") // 3 символа на ячейку (cell=2 + 1 gap)
		}
	}

	fmt.Println(monthLine.String())

	// 7 строк — по одной на день недели
	for row := range 7 {
		fmt.Printf("%s  ", dayLabels[row])

		for w := range weeks {
			day := startMonday.AddDate(0, 0, w*7+row)
			if day.After(today) {
				fmt.Print("   ") // будущие дни — пусто

				continue
			}

			n := counts[day.Format("2006-01-02")]
			style := graphLevels[level(n)]
			fmt.Print(style.Render(cell) + " ")
		}

		fmt.Println()
	}
}

// BarChart prints a horizontal bar chart of string→count data to stdout.
// Used for tag/repo breakdowns in stats.
func BarChart(data map[string]int) {
	_ = data
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

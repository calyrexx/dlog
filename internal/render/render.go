package render

import (
	"github.com/calyrexx/dlog/internal/entities"
)

// Table prints entries as a formatted table to stdout.
// Columns: time, tag, text, repo, branch, duration.
func Table(entries []entities.Entry) {
	// TODO: use text/tabwriter or a table library for aligned output
}

// ContributionGraph prints a GitHub-style contribution heatmap to stdout.
// Each cell represents one day; intensity reflects entry count.
func ContributionGraph(entries []entities.Entry) {
	// TODO: bucket entries by day, render 7-row × N-week grid with ANSI colours
}

// BarChart prints a horizontal bar chart of string→count data to stdout.
// Used for tag/repo breakdowns in stats.
func BarChart(data map[string]int) {
	// TODO: sort keys, determine max value, draw proportional bars with ANSI colour
}

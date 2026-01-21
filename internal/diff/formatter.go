package diff

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("251")).
		Background(lipgloss.Color("238")).
		Padding(0, 1).
		Width(100)
}

var (
	additionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("35"))  // Soft green (moon theme)
	deletionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("167")) // Soft red (moon theme)
	hunkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Soft blue (moon theme)
)

// FormatDiff applies ANSI color formatting to a git diff string
func FormatDiff(width int, diff string) string {
	if diff == "" {
		return ""
	}

	lines := strings.Split(diff, "\n")
	var formatted []string

	for _, line := range lines {
		var styledLine string

		switch {
		case strings.HasPrefix(line, "+"):
			styledLine = additionStyle.Render(line)
		case strings.HasPrefix(line, "-"):
			styledLine = deletionStyle.Render(line)
		case strings.HasPrefix(line, "@@"):
			styledLine = hunkStyle.Render(line)
		case strings.HasPrefix(line, "diff ") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++"):
			styledLine = headerStyle().Width(width).Render(line)
		default:
			styledLine = line
		}

		formatted = append(formatted, styledLine)
	}

	return strings.Join(formatted, "\n")
}

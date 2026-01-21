package diff

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	additionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // Green
	deletionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))   // Red
	hunkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))   // Cyan
	headerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true) // Yellow/Bold
)

// FormatDiff applies ANSI color formatting to a git diff string
func FormatDiff(diff string) string {
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
			styledLine = headerStyle.Render(line)
		default:
			styledLine = line
		}

		formatted = append(formatted, styledLine)
	}

	return strings.Join(formatted, "\n")
}

package diff

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/samverrall/review-ui/internal/ui/color"
)

var (
	additionStyle  = lipgloss.NewStyle().Foreground(color.MoonGreen)  // Soft green (moon theme)
	deletionStyle  = lipgloss.NewStyle().Foreground(color.MoonRed)    // Soft red (moon theme)
	hunkStyle      = lipgloss.NewStyle().Foreground(color.MoonBlue)   // Soft blue (moon theme)
	headerStyle = lipgloss.NewStyle().Foreground(color.MoonPurple) // Muted purple (moon theme)
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
			styledLine = headerStyle.Width(width).Render(line)
		default:
			styledLine = line
		}

		formatted = append(formatted, styledLine)
	}

	return strings.Join(formatted, "\n")
}

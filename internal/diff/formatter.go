package diff

import (
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/samverrall/review-ui/internal/syntax"
	"github.com/samverrall/review-ui/internal/ui/color"
)

var (
	additionStyle = lipgloss.NewStyle().Foreground(color.MoonGreen)  // Soft green (moon theme)
	deletionStyle = lipgloss.NewStyle().Foreground(color.MoonRed)    // Soft red (moon theme)
	hunkStyle     = lipgloss.NewStyle().Foreground(color.MoonBlue)   // Soft blue (moon theme)
	headerStyle   = lipgloss.NewStyle().Foreground(color.MoonPurple) // Muted purple (moon theme)
)

var (
	// Regex to match diff headers and extract filenames
	diffHeaderRegex = regexp.MustCompile(`^diff --git a/(.+) b/(.+)$`)
	fileHeaderRegex = regexp.MustCompile(`^(---|\+\+\+) (.+)$`)
)

var (
	highlighter     *syntax.Highlighter
	highlighterOnce sync.Once
)

// getHighlighter returns the singleton highlighter instance
func getHighlighter() *syntax.Highlighter {
	highlighterOnce.Do(func() {
		highlighter = syntax.NewHighlighter()
	})
	return highlighter
}

// FormatDiff applies ANSI color formatting to a git diff string with syntax highlighting
func FormatDiff(width int, diff string, logger *slog.Logger) string {
	if diff == "" {
		return ""
	}

	h := getHighlighter()
	lines := strings.Split(diff, "\n")
	var formatted []string

	currentFile := "" // Track the current file for syntax highlighting

	for _, line := range lines {
		var styledLine string

		switch {
		case strings.HasPrefix(line, "+"):
			// Addition line - apply syntax highlighting to code, keep + green
			code := strings.TrimPrefix(line, "+")
			if currentFile != "" && code != "" {
				highlighted, err := h.Highlight(currentFile, code)
				if err == nil {
					// Green + symbol followed by syntax-highlighted code
					styledLine = additionStyle.Render("+") + highlighted
				} else {
					logger.Debug("syntax highlighting failed for addition",
						"file", currentFile,
						"error", err)
					styledLine = additionStyle.Render(line)
				}
			} else {
				styledLine = additionStyle.Render(line)
			}
		case strings.HasPrefix(line, "-"):
			// Deletion line - apply syntax highlighting to code, keep - red
			code := strings.TrimPrefix(line, "-")
			if currentFile != "" && code != "" {
				highlighted, err := h.Highlight(currentFile, code)
				if err == nil {
					// Red - symbol followed by syntax-highlighted code
					styledLine = deletionStyle.Render("-") + highlighted
				} else {
					logger.Debug("syntax highlighting failed for deletion",
						"file", currentFile,
						"error", err)
					styledLine = deletionStyle.Render(line)
				}
			} else {
				styledLine = deletionStyle.Render(line)
			}
		case strings.HasPrefix(line, "@@"):
			// Hunk header - keep existing styling
			styledLine = hunkStyle.Render(line)
		case strings.HasPrefix(line, "diff "):
			// Diff header - extract filename for syntax highlighting
			styledLine = headerStyle.Width(width).Render(line)
			if matches := diffHeaderRegex.FindStringSubmatch(line); len(matches) >= 3 {
				// Use the "b/" version (new file) for syntax highlighting
				currentFile = matches[2]
			}
		case strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++"):
			// File header - keep existing styling
			styledLine = headerStyle.Width(width).Render(line)
			if matches := fileHeaderRegex.FindStringSubmatch(line); len(matches) >= 3 {
				// Update current file if this is a +++ line (new file)
				if strings.HasPrefix(line, "+++") && matches[2] != "/dev/null" {
					currentFile = strings.TrimPrefix(matches[2], "b/")
				}
			}
		default:
			// Context line - apply syntax highlighting without tint
			if currentFile != "" {
				highlighted, err := h.Highlight(currentFile, line)
				if err == nil {
					styledLine = highlighted
				} else {
					// Log syntax highlighting failures for debugging
					logger.Debug("syntax highlighting failed",
						"file", currentFile,
						"line", line,
						"error", err)
					styledLine = line
				}
			} else {
				styledLine = line
			}
		}

		formatted = append(formatted, styledLine)
	}

	return strings.Join(formatted, "\n")
}

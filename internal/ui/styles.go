package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/samverrall/review-ui/internal/ui/color"
)

var (
	// Header style for the top bar showing file information - more prominent
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(color.SubtleText).
			Background(color.DarkBg).
			Border(lipgloss.ThickBorder(), false, false, true, false).
			BorderForeground(color.MoonBlue).
			Padding(1, 3).
			Margin(0, 0, 2, 0).
			Align(lipgloss.Center)

	// Footer style for the help text at the bottom
	footerStyle = lipgloss.NewStyle().
			Foreground(color.SubtleText).
			Background(color.DarkBg).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(color.MoonDarkGray).
			Padding(0, 2).
			Margin(2, 0, 0, 0).
			Align(lipgloss.Center)

	// Error style for error messages
	errorStyle = lipgloss.NewStyle().
			Foreground(color.MoonRed).
			Bold(true).
			Background(color.DarkBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(color.MoonRed).
			Padding(1, 2).
			Margin(1, 0)

	// Info style for informational messages
	infoStyle = lipgloss.NewStyle().
			Foreground(color.MoonYellow).
			Background(color.DarkBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(color.MoonYellow).
			Padding(1, 2).
			Margin(1, 0)

	// Cursor line style for highlighting the current line
	cursorLineStyle = lipgloss.NewStyle().
			Background(color.CursorLineBg).
			Foreground(color.TextColor)

	// Selection style for highlighting selected lines
	selectionStyle = lipgloss.NewStyle().
			Background(color.MoonPurple).
			Foreground(color.DarkBg).
			Bold(true)

	// Comment style for displaying comments
	commentStyle = lipgloss.NewStyle().
			Foreground(color.MoonBlue).
			Background(color.AccentBg).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(color.MoonBlue).
			Padding(0, 1).
			Margin(0, 2)

	// Comment input style for the input box
	commentInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(color.MoonPurple).
				Foreground(color.TextColor).
				Background(color.DarkBg).
				Padding(1, 2).
				Margin(1, 0).
				Width(80)

	// Status style for success/error messages
	statusStyle = lipgloss.NewStyle().
			Foreground(color.MoonGreen).
			Background(color.DarkBg).
			Bold(true).
			Padding(0, 2).
			Margin(0, 0, 1, 0)

	// File list styles
	fileListItemStyle = lipgloss.NewStyle().
				Foreground(color.TextColor).
				Padding(0, 1)

	fileListSelectedStyle = lipgloss.NewStyle().
				Foreground(color.DarkBg).
				Background(color.MoonBlue).
				Bold(true).
				Padding(0, 1)

	// Modal container for centered content
	modalContainer = lipgloss.NewStyle().
			Padding(1, 4).
			Margin(0, 2)
)

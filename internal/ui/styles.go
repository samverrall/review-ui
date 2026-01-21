package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Tailwind Moon theme color palette
	moonBlue      = lipgloss.Color("39")   // Soft blue
	moonPurple    = lipgloss.Color("63")   // Muted purple
	moonGreen     = lipgloss.Color("35")   // Soft green
	moonYellow    = lipgloss.Color("221")  // Muted yellow
	moonRed       = lipgloss.Color("167")  // Soft red
	moonGray      = lipgloss.Color("244")  // Light gray
	moonDarkGray  = lipgloss.Color("238")  // Dark gray
	moonDarkerGray = lipgloss.Color("235") // Darker gray

	darkBg        = lipgloss.Color("234")  // Very dark background (like VS Code moon)
	textColor     = lipgloss.Color("251")  // Light text
	subtleText    = lipgloss.Color("244")  // Muted text
	accentBg      = lipgloss.Color("236")  // Accent background

	// Header style for the top bar showing file information - more prominent
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(subtleText).
		Background(darkBg).
		Border(lipgloss.ThickBorder(), false, false, true, false).
		BorderForeground(moonBlue).
		Padding(1, 3).
		Margin(0, 0, 2, 0).
		Align(lipgloss.Center)

	// Footer style for the help text at the bottom
	footerStyle = lipgloss.NewStyle().
		Foreground(subtleText).
		Background(darkBg).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(moonDarkGray).
		Padding(0, 2).
		Margin(2, 0, 0, 0).
		Align(lipgloss.Center)

	// Error style for error messages
	errorStyle = lipgloss.NewStyle().
		Foreground(moonRed).
		Bold(true).
		Background(darkBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(moonRed).
		Padding(1, 2).
		Margin(1, 0)

	// Info style for informational messages
	infoStyle = lipgloss.NewStyle().
		Foreground(moonYellow).
		Background(darkBg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(moonYellow).
		Padding(1, 2).
		Margin(1, 0)

	// Cursor line style for highlighting the current line
	cursorLineStyle = lipgloss.NewStyle().
		Background(accentBg).
		Foreground(textColor)

	// Selection style for highlighting selected lines
	selectionStyle = lipgloss.NewStyle().
		Background(moonPurple).
		Foreground(darkBg).
		Bold(true)

	// Comment style for displaying comments
	commentStyle = lipgloss.NewStyle().
		Foreground(moonBlue).
		Background(accentBg).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(moonBlue).
		Padding(0, 1).
		Margin(0, 2)

	// Comment input style for the input box
	commentInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(moonPurple).
		Foreground(textColor).
		Background(darkBg).
		Padding(1, 2).
		Margin(1, 0).
		Width(80)

	// Status style for success/error messages
	statusStyle = lipgloss.NewStyle().
		Foreground(moonGreen).
		Background(darkBg).
		Bold(true).
		Padding(0, 2).
		Margin(0, 0, 1, 0)

	// File list styles
	fileListItemStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Padding(0, 1)

	fileListSelectedStyle = lipgloss.NewStyle().
		Foreground(darkBg).
		Background(moonBlue).
		Bold(true).
		Padding(0, 1)

	// Modal container for centered content
	modalContainer = lipgloss.NewStyle().
		Padding(1, 4).
		Margin(0, 2)
)

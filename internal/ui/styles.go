package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Header style for the top bar showing file information
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4")).
		Padding(0, 1).
		Width(100)

	// Footer style for the help text at the bottom
	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Padding(0, 1)

	// Error style for error messages
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)

	// Info style for informational messages
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")).
		Padding(1)

	// Cursor line style for highlighting the current line
	cursorLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("237"))

	// Selection style for highlighting selected lines
	selectionStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("229"))

	// Comment style for displaying comments
	commentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Italic(true).
		Padding(0, 2)

	// Comment input style for the input box
	commentInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Padding(0, 1).
		Margin(1, 0)

	// Status style for success/error messages
	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Padding(0, 1)
)

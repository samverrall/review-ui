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
)

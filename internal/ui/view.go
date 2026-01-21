package ui

import (
	"fmt"
	"strings"
)

// View renders the current state of the model
func (m model) View() string {
	// Handle error state
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err))
	}

	// Handle no changes state
	if len(m.changedFiles) == 0 {
		return infoStyle.Render("No unstaged changes found.\n\nPress q to quit.")
	}

	// Handle not ready state (terminal size not yet known)
	if !m.ready {
		return "Initializing..."
	}

	// Build main view
	var b strings.Builder

	// Header: File counter and name
	currentFile := m.changedFiles[m.currentIndex]
	headerText := fmt.Sprintf("File %d/%d: %s", m.currentIndex+1, len(m.changedFiles), currentFile)
	header := headerStyle.Render(headerText)
	if m.width > 0 {
		header = headerStyle.Width(m.width).Render(headerText)
	}
	b.WriteString(header)
	b.WriteString("\n")

	// Viewport: Diff content
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Footer: Help text
	helpText := "n: next | p: prev | ↑↓: scroll | q: quit"
	footer := footerStyle.Render(helpText)
	b.WriteString(footer)

	return b.String()
}
